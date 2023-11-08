package lib

import (
	"github.com/skratchdot/open-golang/open"
)

// var CB_URI = "http://localhost:8088/callback"

func SignUp() {
	open.Start("https://pdn-vehpksfu665ae7k5jewmycb4fxqircam.login.aws.us.pangea.cloud")

}

// func flowHandlePasswordPhase(ctx context.Context, client *authn.AuthN, flow_id, password string) *authn.FlowUpdateResult {
// 	fmt.Println("Handling password phase...")
// 	resp, err := client.Flow.Update(ctx, authn.FlowUpdateRequest{
// 		FlowID: flow_id,
// 		Choice: authn.FCPassword,
// 		Data: authn.FlowUpdateDataPassword{
// 			Password: password,
// 		},
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return resp.Result
// }

// func flowHandleProfilePhase(ctx context.Context, client *authn.AuthN, flow_id string, profile *authn.ProfileData) *authn.FlowUpdateResult {
// 	fmt.Println("Handling profile phase...")
// 	resp, err := client.Flow.Update(ctx, authn.FlowUpdateRequest{
// 		FlowID: flow_id,
// 		Choice: authn.FCProfile,
// 		Data: authn.FlowUpdateDataProfile{
// 			Profile: *profile,
// 		},
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return resp.Result
// }

// func flowHandleAgreementsPhase(ctx context.Context, client *authn.AuthN, flow_id string, result *authn.FlowUpdateResult) *authn.FlowUpdateResult {
// 	// Iterate over flow_choices in response.result
// 	fmt.Println("Handling agreements phase...")
// 	agreed := []string{}
// 	for _, flowChoice := range result.FlowChoices {
// 		// Check if the choice is AGREEMENTS
// 		if flowChoice.Choice == string(authn.FCAgreements) {
// 			// Assuming flowChoice.Data["agreements"] is a map[string]interface{}
// 			agreements, ok := flowChoice.Data["agreements"].(map[string]interface{})
// 			if ok {
// 				// Iterate over agreements and append the "id" values to agreed slice
// 				for _, v := range agreements {
// 					agreement, ok := v.(map[string]interface{})
// 					if ok {
// 						id, ok := agreement["id"].(string)
// 						if ok {
// 							agreed = append(agreed, id)
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}

// 	resp, err := client.Flow.Update(ctx, authn.FlowUpdateRequest{
// 		FlowID: flow_id,
// 		Choice: authn.FCAgreements,
// 		Data: authn.FlowUpdateDataAgreements{
// 			Agreed: agreed,
// 		},
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return resp.Result
// }

// func choiceIsAvailable(choices []authn.FlowChoiceItem, choice string) bool {
// 	for _, fc := range choices {
// 		if fc.Choice == choice {
// 			return true

// 		}
// 	}
// 	return false
// }

// func CreateAndLogin(first_name string, last_name string, email string, password string) *authn.FlowCompleteResult {
// 	r := mux.NewRouter()

// 	r.HandleFunc("/callback", controller.Callback).Methods("GET")

// 	l, err := net.Listen("tcp", "localhost:8088")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	go func() {
// 		if err := http.Serve(l, r); err != nil {
// 			log.Fatal(err)
// 		}
// 	}()

// 	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancelFn()

// 	client := Init()

// 	profile := &authn.ProfileData{
// 		"first_name": first_name,
// 		"last_name":  last_name,
// 	}

// 	fmt.Println("Flow starting...")
// 	fsresp, err := client.Flow.Start(ctx,
// 		authn.FlowStartRequest{
// 			Email:     email,
// 			FlowTypes: []authn.FlowType{authn.FTsignup, authn.FTsignin},
// 			CBURI:     CB_URI,
// 		})
// 	// fsresp, err := authn.Flow.Start(ctx,
// 	// 	authn.FlowStartRequest{
// 	// 		Email:     email,
// 	// 		FlowTypes: []authn.FlowType{authn.FTsignup, authn.FTsignin},
// 	// 		CBURI:     CB_URI,
// 	// 	})

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	flowID := fsresp.Result.FlowID
// 	var result *authn.FlowUpdateResult
// 	flowPhase := "initial"
// 	choices := fsresp.Result.FlowChoices
// 	fmt.Println(fsresp.Result.FlowChoices)

// 	for flowPhase != "phase_completed" {
// 		if choiceIsAvailable(choices, string(authn.FCPassword)) {
// 			fmt.Printf("true")
// 			result = flowHandlePasswordPhase(ctx, client, flowID, password)
// 		} else if choiceIsAvailable(choices, string(authn.FCProfile)) {
// 			fmt.Printf("true")
// 			result = flowHandleProfilePhase(ctx, client, flowID, profile)
// 		} else if choiceIsAvailable(choices, string(authn.FCAgreements)) {
// 			fmt.Printf("true")
// 			result = flowHandleAgreementsPhase(ctx, client, flowID, result)
// 		} else {
// 			if result != nil {
// 				fmt.Printf("Phase %s not handled", result.FlowPhase)
// 			} else {
// 				fmt.Printf("Phase not handled, result is nil")
// 			}
// 		}
// 		if result != nil {
// 			fmt.Printf("true11")
// 			flowPhase = result.FlowPhase
// 			choices = result.FlowChoices
// 		} else {
// 			fmt.Printf("true1")
// 			break
// 		}
// 	}

// 	fcresp, err := client.Flow.Complete(ctx,
// 		authn.FlowCompleteRequest{
// 			FlowID: flowID,
// 		})

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	user, err := user.Current()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	filePath := filepath.Join(user.HomeDir, "Securelee/tokens.json")

// 	data, err := json.Marshal(fcresp.Result)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	err = ioutil.WriteFile(filePath, data, 0644)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return fcresp.Result
// }

func Logout() {

}
