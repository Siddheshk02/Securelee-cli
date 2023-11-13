package lib

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/vault"
)

func Create(ID string) error {
	var choice int
	fmt.Print("\033[33m", "\n Select any one option: \n", "\033[0m")
	fmt.Println("\033[33m", "\n > 1. Store a Secret Message", "\033[0m")
	fmt.Println("\033[33m", "> 2. Create or Store a Key", "\033[0m")
	fmt.Println("")
	fmt.Print("\033[36m", " > Enter your choice (for e.g. 1): ", "\033[0m")
	fmt.Scanf("%d", &choice)
	fmt.Println("")
	var str string

	if choice == 1 {
		str = "secrets"
	} else if choice == 2 {
		str = "keys"
	} else {
		fmt.Println("\033[31m", "\n > Invalid Choice Entered!!, Please try again", "\033[0m")
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
		fmt.Print("\033[33m", " > Enter the Secret Value ", "\033[0m")
		fmt.Print("\033[35m", "(no spaces in-between) {Minimum Size : 1, Maximum Size: 10240} : ", "\033[0m")
		fmt.Scan(&secretvalue)
		fmt.Println("")
		fmt.Print("\033[33m", " > Enter the Name for the Secret ", "\033[0m")
		fmt.Print("\033[35m", "(no spaces in-between): ", "\033[0m")
		fmt.Scan(&secretname)

		fmt.Println("")
		// time.Sleep(15 * time.Second)

		input := &vault.SecretStoreRequest{
			CommonStoreRequest: vault.CommonStoreRequest{
				Name:   secretname,
				Folder: folder,
			},
			Secret: secretvalue,
		}
		rStore, err := client.SecretStore(ctx, input)
		if err != nil {
			return err
		}

		if *rStore.Status != "Success" {
			fmt.Println("\033[31m", " > A Secret with same name exists.", "\033[0m")
			os.Exit(0)
		}

		fmt.Println("\033[36m", " > Secret Successfully Stored!!", "\033[0m")

		return nil

	} else {

		res, err, ch := StoreKey(folder)
		if res == "" && err != nil {
			return err
		} else if res != "" && err != nil {
			fmt.Println("\033[31m", "\n > ", res, "\033[0m")
			fmt.Println("")
			return nil
		} else if res != "" && err == nil {
			if res == "Success" {
				if ch == 1 {
					fmt.Println("\033[36m", " > Key Stored Successfully!", "\033[0m")
					fmt.Println("")
					return nil
				} else {
					fmt.Println("\033[36m", " > Key Genarated Successfully!", "\033[0m")
					fmt.Println("")
					return nil
				}
			} else {
				fmt.Println("\033[36m", " > The Public Key : ", res, "\033[0m")
				return nil
			}
		}
		return nil
	}

}

func StoreKey(folder string) (string, error, int) {
	var choice int
	fmt.Print("\033[33m", "\n Select any one option: \n", "\033[0m")
	fmt.Println("\033[33m", "\n > 1. Import a Key", "\033[0m")
	fmt.Println("\033[33m", "> 2. Generate a Key", "\033[0m")
	// fmt.Println("")
	fmt.Print("\033[36m", "\n > Enter your choice (for e.g. 1): ", "\033[0m")
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
		fmt.Println("\033[33m", "\n Select the Key type: ", "\033[0m")
		fmt.Println("\033[33m", "\n > 1. Asymmetric Key", "\033[0m")
		fmt.Println("\033[33m", "> 2. Symmetric Key", "\033[0m")
		fmt.Println("")
		fmt.Print("\033[36m", " > Enter your choice (for e.g. 1): ", "\033[0m")
		// fmt.Scanf("%d", &choice2)
		fmt.Scan(&choice2)
		fmt.Println("")

		if choice2 == 1 {
			var keyname, keypurpose, publickey, privatekey string
			fmt.Print("\033[33m", " > Enter the Name for the Key ", "\033[0m")
			fmt.Print("\033[35m", "(no spaces in-between): ", "\033[0m")
			fmt.Scan(&keyname)
			fmt.Println("")
			fmt.Print("\033[33m", " > Enter the Purpose for the Key ", "\033[0m")
			fmt.Print("\033[35m", "( for e.g. signing, Encryption): ", "\033[0m")
			fmt.Scan(&keypurpose)
			fmt.Print("\n")
			fmt.Print("\033[33m", " > Enter the Public Key ", "\033[0m")
			fmt.Print("\033[35m", "(in PEM format): ", "\033[0m")
			fmt.Scan(&publickey)
			fmt.Println("")
			fmt.Print("\033[33m", " > Enter the Private Key ", "\033[0m")
			fmt.Print("\033[35m", "(in PEM format): ", "\033[0m")
			fmt.Scan(&privatekey)
			fmt.Println("")

			rStore, err := client.AsymmetricStore(ctx,
				&vault.AsymmetricStoreRequest{
					CommonStoreRequest: vault.CommonStoreRequest{
						Name:          keyname,
						Folder:        folder,
						RotationState: vault.IVSdeactivated,
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
				err1 := json.Unmarshal(match, &apiError)
				if err1 != nil {
					return "", err1, choice
				}

				parts := strings.Split(apiError.Summary, ": ")

				content := strings.TrimSpace(parts[1])

				return content, err, choice
			} else if err == nil && rStore != nil {
				if *rStore.Status == "Success" {
					return string(rStore.Result.PublicKey), nil, choice
				}
			}

		} else if choice2 == 2 {
			var keyname, key string
			fmt.Print("\033[33m", " > Enter the Name for the Key ", "\033[0m")
			fmt.Print("\033[35m", "(no spaces in-between): ", "\033[0m")
			fmt.Scan(&keyname)
			fmt.Println("")
			fmt.Print("\033[33m", " > Enter the Key ", "\033[0m")
			fmt.Print("\033[35m", "(base64): ", "\033[0m")
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
			fmt.Println("\033[31m", "\n > Invalid Choice Entered!!, Please try again", "\033[0m")
			fmt.Println("")
			os.Exit(0)
		}

	} else if choice == 2 {
		ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancelFn()
		client := InitVault()
		var choice2 int
		fmt.Println("\033[33m", "\n Select the Key type: ", "\033[0m")
		fmt.Println("\033[33m", "\n > 1. Asymmetric Key", "\033[0m")
		fmt.Println("\033[33m", "> 2. Symmetric Key", "\033[0m")
		fmt.Println("")
		fmt.Print("\033[36m", " > Enter your choice (for e.g. 1): ", "\033[0m")
		// fmt.Scanf("%d", &choice2)
		fmt.Scan(&choice2)
		fmt.Println("")

		if choice2 == 1 {
			var keyname, keypurpose string
			fmt.Print("\033[33m", " > Enter the Name for the Key ", "\033[0m")
			fmt.Print("\033[35m", "(no spaces in-between): ", "\033[0m")
			fmt.Scan(&keyname)
			fmt.Println("")
			fmt.Print("\033[33m", " > Enter the Purpose for the Key ", "\033[0m")
			fmt.Print("\033[35m", "( for e.g. signing, Encryption): ", "\033[0m")
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
			fmt.Print("\033[33m", " > Enter the Name for the Key ", "\033[0m")
			fmt.Print("\033[35m", "(no spaces in-between): ", "\033[0m")
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
			fmt.Println("\033[31m", "\n > Invalid Choice Entered!!, Please try again", "\033[0m")
			fmt.Println("")
			os.Exit(0)
		}

	} else {
		fmt.Println("\033[31m", "\n > Invalid Choice Entered!!, Please try again", "\033[0m")
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
		fmt.Println("\033[33m", "List of Secrets : ", "\033[0m")

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

		str := " "
		fmt.Fprintf(w, " +-------+-------------------------------------------+----------------------------+--------+-----------------+\n")
		fmt.Fprintf(w, " | Index | ID%-40s| Name %-21s | Type   | State %-9s |", str, str, str)
		fmt.Fprintf(w, "\n +-------+-------------------------------------------+----------------------------+--------+-----------------+\n")

		for i := 0; i < resp.Result.Count; i++ {
			fmt.Fprintf(w, " | %-5d | %-40s  | %-25s  | %-4s | %-15s |\n", i+1, lists[i].ID, lists[i].Name, lists[i].Type, lists[i].CurrentVersion.State)
		}
		fmt.Fprintf(w, " +-------+-------------------------------------------+----------------------------+--------+-----------------+\n")

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
		fmt.Println("\033[33m", "List of Keys : ", "\033[0m")

		lists := resp.Result.Items
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

		str := " "
		fmt.Fprintf(w, " +-------+---------------------------------------+----------------------------+----------------+---------------+---------------------------+-----------------+\n")
		fmt.Fprintf(w, " | Index | ID%-36s| Name %-21s | Type  %-8s | Purpose %-5s | Algorithm %-15s | State %-9s |", str, str, str, str, str, str)
		fmt.Fprintf(w, "\n +-------+---------------------------------------+----------------------------+----------------+---------------+---------------------------+-----------------+\n")

		for i := 0; i < resp.Result.Count; i++ {
			fmt.Fprintf(w, " | %-5d | %-36s  | %-25s  | %-15s| %-13s | %-25s | %-15s |\n", i+1, lists[i].ID, lists[i].Name, lists[i].Type, lists[i].Purpose, lists[i].Algorithm, lists[i].CurrentVersion.State)
		}
		fmt.Fprintf(w, " +-------+---------------------------------------+----------------------------+----------------+---------------+---------------------------+-----------------+\n")

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
			fmt.Println("\033[31m", "\n > ID doesn't exists. Please try again.\n", "\033[0m")
			os.Exit(0)
		}
	}

	if *dresp.Status == "Success" {
		fmt.Println("\033[36m", "\n > Item deleted Successfully!\n", "\033[0m")
		os.Exit(0)
	}

	return nil
}

func Update(id string, req string) error {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := InitVault()

	var name string
	fmt.Print("\033[33m", "\n > Enter the Name for the ", req, " : ", "\033[0m")
	fmt.Scan(&name)
	fmt.Println("")

	var state string
	var choice int

	fmt.Println("\033[33m", "\n Select the Item State for the", req, ": ", "\033[0m")
	fmt.Println("\033[33m", "\n > 1. Active", "\033[0m")
	fmt.Println("\033[33m", "> 2. Deactivated", "\033[0m")
	fmt.Println("\033[33m", "> 3. Suspended", "\033[0m")
	fmt.Println("\033[33m", "> 4. Compromised", "\033[0m")
	fmt.Println("\033[33m", "> 5. Destroyed", "\033[0m")
	// fmt.Println("> 1. Enable")
	// fmt.Println("> 2. Inherited")
	fmt.Println("")
	fmt.Print("\033[36m", " > Enter your choice (for e.g. 1): ", "\033[0m")
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
	default:
		fmt.Println("\033[31m", "Invalid Choice! Please try again.", "\033[0m")
		os.Exit(0)
	}

	rUpdate, err1 := client.Update(ctx,
		&vault.UpdateRequest{
			ID:   id,
			Name: name,
		},
	)

	if err1 != nil {
		re := regexp.MustCompile(`\{[^{}]*\}`)
		match := re.Find([]byte(err1.Error()))

		if match == nil {
			return errors.New("No JSON data found in the error message")
		}

		var apiError APIError
		err := json.Unmarshal(match, &apiError)
		if err != nil {
			return err
		}

		if apiError.Status == "VaultItemNotFound" {
			fmt.Println("\033[31m", "\n > ID doesn't exists. Please try again.", "\033[0m")
			fmt.Println("")
			os.Exit(0)
		}
	}

	input := &vault.StateChangeRequest{
		ID:    id,
		State: vault.ItemVersionState(state),
	}

	scr, err2 := client.StateChange(ctx, input)

	if err2 != nil {
		re := regexp.MustCompile(`\{[^{}]*\}`)
		match := re.Find([]byte(err2.Error()))

		if match == nil {
			return errors.New("No JSON data found in the error message")
		}

		var apiError APIError
		err := json.Unmarshal(match, &apiError)
		if err != nil {
			return err
		}

		if apiError.Status != "Success" {
			fmt.Println("\033[35m", "\n > Item Name Updated!. State can't be Updated.", "\033[0m")
			fmt.Println("")
			os.Exit(0)
		}
	}

	if *rUpdate.Status == "Success" && *scr.Status == "Success" {
		fmt.Println("\033[36m", "\n > Item Updated Successfully!", "\033[0m")
		fmt.Println("")
		os.Exit(0)
	}

	return nil
}
