package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	swf "github.com/aws/aws-sdk-go/service/swf"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	QueueName  string
	WorkflowID string
)

const (
	DEFAULT_TASK_LIST = "defaultTaskList"
)

func init() {
	workerCmd.Flags().StringVar(&QueueName, "queue", "Main", "Queue")
	workerCmd.Flags().StringVar(&WorkflowID, "id", "", "WorkflowID")
}

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Run worker",
	Run: func(cmd *cobra.Command, args []string) {

		err := godotenv.Load()
		if err != nil {
			fmt.Printf("Error loading .env file")
			os.Exit(1)
		}

		mySession := session.Must(session.NewSession())
		client := swf.New(mySession)

		if WorkflowID == "" {
			fmt.Println("Registering domain/namespace")
			tasks := &swf.TaskList{}
			tasks.SetName(DEFAULT_TASK_LIST)
			input := &swf.RegisterWorkflowTypeInput{
				DefaultTaskList: tasks,
			}
			input.SetName("aws-swf-er-poc")
			input.SetDomain("default")
			input.SetVersion("1")
			input.SetDefaultChildPolicy("TERMINATE")
			input.SetDefaultExecutionStartToCloseTimeout("3600")
			input.SetDefaultTaskStartToCloseTimeout("3600")
			result, err := client.RegisterWorkflowType(input)
			if err != nil {
				log.Fatalf("Error while registering workflow type:%s", err.Error())
			}
			fmt.Println("Results:" + result.GoString())
			os.Exit(0)
			return
		}

		for {
			pdi := &swf.PollForDecisionTaskInput{}
			pdi.SetDomain("default")
			pdi.SetIdentity(WorkflowID)
			pdo, err := client.PollForDecisionTask(pdi)
			if err != nil {
				log.Fatalf("Error while polling for decision:%s", err.Error())
			}
			if pdo.TaskToken == nil {
				continue
			}
			processEvents(
				client,
				pdo.WorkflowExecution.WorkflowId,
				pdo.TaskToken,
				pdo.Events)
		}

	},
}

func processEvents(client *swf.SWF, workflowId *string, taskToken *string, events []*swf.HistoryEvent) {
	for _, event := range events {
		switch *event.EventType {
		case "":
			break
		}
		decisions := []*swf.Decision{}

		decision := &swf.Decision{}
		decision.SetDecisionType("ScheduleActivityTask")

		attribute := &swf.ScheduleActivityTaskDecisionAttributes{}
		attribute.SetActivityId("activity1")

		activityType := &swf.ActivityType{}
		activityType.SetName("name1")
		activityType.SetVersion("1")
		attribute.SetActivityType(activityType)

		decision.SetScheduleActivityTaskDecisionAttributes(attribute)

		decisions = append(decisions, decision)

		input := &swf.RespondDecisionTaskCompletedInput{}
		input.SetTaskList(&swf.TaskList{})
		input.SetTaskToken(*taskToken)
		input.SetDecisions(decisions)
		client.RespondDecisionTaskCompleted(input)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() error {
	// workerCmd.Use = appName

	// Silence Cobra's internal handling of command usage help text.
	// Note, the help text is still displayed if any command arg or
	// flag validation fails.
	workerCmd.SilenceUsage = true

	// Silence Cobra's internal handling of error messaging
	// since we have a custom error handler in main.go
	workerCmd.SilenceErrors = true

	err := workerCmd.Execute()
	return err
}
