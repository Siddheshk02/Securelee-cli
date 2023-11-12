package lib

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/vault"
)

func Create(ID string) error {
	var choice int
	fmt.Println("\n Select any one option: ")
	fmt.Println("\n> 1. Store a Secret Message")
	fmt.Println("> 2. Create or Store a Key")
	fmt.Println("")
	fmt.Print("> Enter your choice (for e.g. 1): ")
	fmt.Scanf("%d", &choice)
	fmt.Println("")
	var str string

	if choice == 1 {
		str = "secrets"
	} else if choice == 2 {
		str = "keys"
	} else {
		fmt.Println("> Invalid Choice Entered!!, Please try again")
		fmt.Println("")
		os.Exit(0)
	}

	foldername := "/" + ID

	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := InitVault()

	input := &vault.FolderCreateRequest{
		Name:   str,
		Folder: foldername,
	}

	_, _ = client.FolderCreate(ctx, input)
	// if err != nil {
	// 	continue
	// }
	folder := foldername + "/" + str

	if choice == 1 {
		var secretvalue, secretname string
		fmt.Print("> Enter the Secret Value {Minimum Size : 1, Maximum Size: 10240} : ")
		fmt.Scan(&secretvalue)
		fmt.Println("")
		fmt.Print("> Enter the Name for the Secret (no spaces in-between): ")
		fmt.Scan(&secretname)
		time.Sleep(7 * time.Second)

		input := &vault.SecretStoreRequest{
			CommonStoreRequest: vault.CommonStoreRequest{
				Name:          secretname,
				Folder:        folder,
				RotationState: vault.IVSactive,
			},
			Secret: secretvalue,
		}
		rStore, err := client.SecretStore(ctx, input)
		if err != nil {
			return err
		}

		if *rStore.Status != "Success" {
			log.Fatal("\n> A Secret with same name exists.")
		}

		fmt.Println("\n> Secret Successfully Stored!!")

		return nil

	} else {

		res, err, ch := StoreKey(folder)
		if res == "" && err != nil {
			return err
		} else if res != "" && err != nil {
			fmt.Println("\n> ", res)
			fmt.Println("")
			return nil
		} else if res != "" && err == nil {
			if res == "Success" {
				if ch == 1 {
					fmt.Println("> Key Stored Successfully!")
					fmt.Println("")
					return nil
				} else {
					fmt.Println("> Key Genarated Successfully!")
					fmt.Println("")
					return nil
				}
			} else {
				fmt.Println("> The Public Key : ", res)
				return nil
			}
		}
		return nil
	}

}

func StoreKey(folder string) (string, error, int) {
	var choice int
	fmt.Println("\n Select any one option: ")
	fmt.Println("\n> 1. Import a Key")
	fmt.Println("> 2. Generate a Key")
	// fmt.Println("")
	fmt.Print("\n> Enter your choice (for e.g. 1): ")
	fmt.Scan(&choice)
	// _, err := fmt.Scanf("%d", &choice)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	fmt.Println("")
	// fmt.Println(choice)

	if choice == 1 {
		ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancelFn()
		client := InitVault()
		var choice2 int
		fmt.Println("\n Select the Key type: ")
		fmt.Println("\n> 1. Asymmetric Key")
		fmt.Println("> 2. Symmetric Key")
		fmt.Println("")
		fmt.Print("> Enter your choice (for e.g. 1): ")
		// fmt.Scanf("%d", &choice2)
		fmt.Scan(&choice2)
		fmt.Println("")

		if choice2 == 1 {
			var keyname, keypurpose, publickey, privatekey string
			fmt.Print("> Enter the Name for the Key (no spaces in-between): ")
			fmt.Scan(&keyname)
			fmt.Println("")
			fmt.Print("> Enter the Purpose for the Key ( for e.g. signing, Encryption): ")
			fmt.Scan(&keypurpose)
			fmt.Println("")
			fmt.Print("> Enter the Public Key (in PEM format): ")
			fmt.Scan(&publickey)
			fmt.Println("")
			fmt.Print("> Enter the Private Key (in PEM format): ")
			fmt.Scan(&privatekey)
			fmt.Println("")

			rStore, err := client.AsymmetricStore(ctx,
				&vault.AsymmetricStoreRequest{
					CommonStoreRequest: vault.CommonStoreRequest{
						Name:          keyname,
						Folder:        folder,
						RotationState: vault.IVSactive,
					},
					Algorithm:  vault.AAed25519,
					Purpose:    vault.KeyPurpose(keypurpose),
					PublicKey:  vault.EncodedPublicKey(publickey),
					PrivateKey: vault.EncodedPrivateKey(privatekey),
				})

			if err != nil && rStore == nil {
				re := regexp.MustCompile(`\{[^{}]*\}`)
				match := re.Find([]byte(err.Error()))

				if match == nil {
					return "", errors.New("No JSON data found in the error message"), choice
				}

				var apiError APIError
				err = json.Unmarshal(match, &apiError)
				if err != nil {
					return "", err, choice
				}

				parts := strings.Split(apiError.Summary, ":")

				content := strings.TrimSpace(parts[1])

				return content, err, choice
			} else if err == nil && rStore != nil {
				if *rStore.Status == "Success" {
					return string(rStore.Result.PublicKey), nil, choice
				}
			}

		} else if choice2 == 2 {
			var keyname, key string
			fmt.Print("> Enter the Name for the Key (no spaces in-between): ")
			fmt.Scan(&keyname)
			fmt.Println("")
			fmt.Print("> Enter the Key (base64): ")
			fmt.Scan(&key)
			fmt.Println("")

			rStore, err := client.SymmetricStore(ctx,
				&vault.SymmetricStoreRequest{
					CommonStoreRequest: vault.CommonStoreRequest{
						Name:          keyname,
						Folder:        folder,
						RotationState: vault.IVSactive,
					},
					Algorithm: vault.SYAaes,
					Purpose:   vault.KPencryption,
					Key:       vault.EncodedSymmetricKey(key),
				})

			if err != nil && rStore == nil {
				re := regexp.MustCompile(`\{[^{}]*\}`)
				match := re.Find([]byte(err.Error()))

				if match == nil {
					return "", errors.New("No JSON data found in the error message"), choice
				}

				var apiError APIError
				err = json.Unmarshal(match, &apiError)
				if err != nil {
					return "", err, choice
				}

				parts := strings.Split(apiError.Summary, ":")

				content := strings.TrimSpace(parts[1])

				return content, err, choice
			} else if err == nil && rStore != nil {
				if *rStore.Status == "Success" {
					return "Success", nil, choice
				}
			}
		} else {
			fmt.Println("> Invalid Choice Entered!!, Please try again")
			fmt.Println("")
			os.Exit(0)
		}

	} else if choice == 2 {
		ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancelFn()
		client := InitVault()
		var choice2 int
		fmt.Println("\n Select the Key type: ")
		fmt.Println("\n> 1. Asymmetric Key")
		fmt.Println("> 2. Symmetric Key")
		fmt.Println("")
		fmt.Print("> Enter your choice (for e.g. 1): ")
		// fmt.Scanf("%d", &choice2)
		fmt.Scan(&choice2)
		fmt.Println("")

		if choice2 == 1 {
			var keyname, keypurpose string
			fmt.Print("> Enter the Name for the Key (no spaces in-between): ")
			fmt.Scan(&keyname)
			fmt.Println("")
			fmt.Print("> Enter the Purpose for the Key ( for e.g. signing, Encryption): ")
			fmt.Scan(&keypurpose)
			fmt.Println("")

			_, err := client.AsymmetricGenerate(ctx,
				&vault.AsymmetricGenerateRequest{
					CommonGenerateRequest: vault.CommonGenerateRequest{
						Name:          keyname,
						Folder:        folder,
						RotationState: vault.IVSactive,
					},
					Algorithm: vault.AArsa2048_pkcs1v15_sha256,
					Purpose:   vault.KeyPurpose(keypurpose),
				})

			if err != nil {
				return "", err, choice
			} else {
				return "Success", nil, choice
			}
		} else if choice2 == 2 {
			var keyname string
			fmt.Print("> Enter the Name for the Key (no spaces in-between): ")
			fmt.Scan(&keyname)
			fmt.Println("")

			_, err := client.SymmetricGenerate(ctx,
				&vault.SymmetricGenerateRequest{
					CommonGenerateRequest: vault.CommonGenerateRequest{
						Name:          keyname,
						Folder:        folder,
						RotationState: vault.IVSactive,
					},
					Algorithm: vault.SYAaes,
					Purpose:   vault.KPencryption,
				})

			if err != nil {
				return "", err, choice
			} else {
				return "Success", nil, choice
			}
		} else {
			fmt.Println("> Invalid Choice Entered!!, Please try again")
			fmt.Println("")
			os.Exit(0)
		}

	} else {
		fmt.Println("> Invalid Choice Entered!!, Please try again")
		fmt.Println("")
		os.Exit(0)
	}

	return "", nil, choice
}

func ListSecrets(UserId string) error {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := InitVault()

	resp, err := client.List(ctx,
		&vault.ListRequest{
			Filter: map[string]string{
				"folder": "/" + UserId + "/secrets",
			},
		})

	if err != nil {
		return err
	} else if *resp.Status == "Success" {
		lists := resp.Result.Items
		fmt.Println("List of Secrets : ")

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

		str := " "
		fmt.Fprintf(w, "+-------+-------------------------------------------+----------------------------+--------+-----------------+\n")
		fmt.Fprintf(w, "| Index | ID%-40s| Name %-21s | Type   | State %-9s |", str, str, str)
		fmt.Fprintf(w, "\n+-------+-------------------------------------------+----------------------------+--------+-----------------+\n")

		for i := 0; i < resp.Result.Count; i++ {
			fmt.Fprintf(w, "| %-5d | %-40s  | %-25s  | %-4s | %-15s |\n", i+1, lists[i].ID, lists[i].Name, lists[i].Type, lists[i].CurrentVersion.State)
		}
		fmt.Fprintf(w, "+-------+-------------------------------------------+----------------------------+--------+-----------------+\n")

		w.Flush()
		return nil
	}

	return nil
}

func ListKeys(UserId string) error {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := InitVault()

	resp, err := client.List(ctx,
		&vault.ListRequest{
			Filter: map[string]string{
				"folder": "/" + UserId + "/keys",
			},
		})

	if err != nil {
		return err
	} else if *resp.Status == "Success" {
		fmt.Println("List of Keys : ")

		lists := resp.Result.Items
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

		str := " "
		fmt.Fprintf(w, "+-------+---------------------------------------+----------------------------+----------------+---------------+---------------------------+-----------------+\n")
		fmt.Fprintf(w, "| Index | ID%-36s| Name %-21s | Type  %-8s | Purpose %-5s | Algorithm %-15s | State %-9s |", str, str, str, str, str, str)
		fmt.Fprintf(w, "\n+-------+---------------------------------------+----------------------------+----------------+---------------+---------------------------+-----------------+\n")

		for i := 0; i < resp.Result.Count; i++ {
			fmt.Fprintf(w, "| %-5d | %-36s  | %-25s  | %-15s| %-13s | %-25s | %-15s |\n", i+1, lists[i].ID, lists[i].Name, lists[i].Type, lists[i].Purpose, lists[i].Algorithm, lists[i].CurrentVersion.State)
		}
		fmt.Fprintf(w, "+-------+---------------------------------------+----------------------------+----------------+---------------+---------------------------+-----------------+\n")

		w.Flush()

		return nil
	}

	return nil
}

func Delete(id string) error {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := InitVault()

	dresp, err := client.Delete(ctx, &vault.DeleteRequest{
		ID: id,
	})

	if err != nil {
		re := regexp.MustCompile(`\{[^{}]*\}`)
		match := re.Find([]byte(err.Error()))

		if match == nil {
			return errors.New("No JSON data found in the error message")
		}

		var apiError APIError
		err = json.Unmarshal(match, &apiError)
		if err != nil {
			return err
		}

		if apiError.Status == "VaultItemNotFound" {
			fmt.Println("\n> ID doesn't exists. Please try again.\n")
			os.Exit(0)
		}
	}

	if *dresp.Status == "Success" {
		fmt.Println("\n> Item deleted Successfully!\n")
		os.Exit(0)
	}

	return nil
}

func Update(id string, req string) error {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := InitVault()

	var name string
	fmt.Print("\n> Enter the Name for the ", req, " : ")
	fmt.Scan(&name)
	fmt.Println("")

	var state string
	var choice int

	fmt.Println("\n Select the Item State for the", req, ": ")
	fmt.Println("\n> 1. Active")
	fmt.Println("> 2. Deactivated")
	fmt.Println("> 3. Suspended")
	fmt.Println("> 4. Compromised")
	fmt.Println("> 5. Destroyed")
	fmt.Println("> 6. Inherited")
	// fmt.Println("> 1. Enable")
	// fmt.Println("> 2. Inherited")
	fmt.Println("")
	fmt.Print("> Enter your choice (for e.g. 1): ")
	fmt.Scan(&choice)
	fmt.Println("")

	switch choice {
	case 1:
		state = "active"
	case 2:
		state = "deactivated"
	case 3:
		state = "suspended"
	case 4:
		state = "compromised"
	case 5:
		state = "destroyed"
	case 6:
		state = "inherited"
	default:
		fmt.Println("Invalid Choice! Please try again.")
		os.Exit(0)
	}

	rUpdate, err := client.Update(ctx,
		&vault.UpdateRequest{
			ID:   id,
			Name: name,
		},
	)

	if err != nil {
		fmt.Println(err)
		re := regexp.MustCompile(`\{[^{}]*\}`)
		match := re.Find([]byte(err.Error()))

		if match == nil {
			return errors.New("No JSON data found in the error message")
		}

		var apiError APIError
		err = json.Unmarshal(match, &apiError)
		if err != nil {
			return err
		}

		if apiError.Status == "VaultItemNotFound" {
			fmt.Println("\n> ID doesn't exists. Please try again.")
			fmt.Println("")
			os.Exit(0)
		}
	}

	input := &vault.StateChangeRequest{
		ID:    id,
		State: vault.ItemVersionState(state),
	}

	scr, err := client.StateChange(ctx, input)

	if *rUpdate.Status == "Success" && *scr.Status == "Success" {
		fmt.Println("\n> Item Updated Successfully!")
		fmt.Println("")
		os.Exit(0)
	}

	return nil
}
