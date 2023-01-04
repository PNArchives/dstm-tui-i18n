package cluster

import (
	"fmt"

	"github.com/PNCommand/dstm/tui/component"
	"github.com/PNCommand/dstm/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func NewCluster(clusterName string) {
	if utils.IsClusterExists(clusterName) {
		fmt.Printf("Cluster %s already exists\n", clusterName)
		return
	}
	token, err := getToken()
	fmt.Println("Token:", token)
	if err != nil {
		fmt.Println(err)
		return
	}
	model, err := generateClusterIniModel()
	if err != nil {
		fmt.Println(err)
		return
	}
	generateClusterFiles(clusterName, token, model)
}

func getToken() (string, error) {
	scanner := component.NewScanner("enter cluster token")
	model, err := tea.NewProgram(scanner).Run()
	if err != nil {
		return "", err
	}
	token := model.(*component.Scanner).GetInput()
	return token, nil
}

func generateClusterIniModel() (*component.IniEditor, error) {
	editor, err := component.NewIniEditor(false, "ini editor")
	if err != nil {
		return nil, err
	}
	model, err := tea.NewProgram(editor).Run()
	if err != nil {
		return nil, err
	}
	return model.(*component.IniEditor), nil
}

func generateClusterFiles(clusterName, token string, model *component.IniEditor) error {
	fmt.Println(clusterName, token)
	return nil
}
