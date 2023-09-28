package data_manager

import "jsondb/internal/files_manager"

// type ProjectsManager struct {
// }

func CreateProject(kwargs map[string]interface{}, MAIN_DIR_PATH string) error {
	return files_manager.CreateDir(MAIN_DIR_PATH, kwargs["name"].(string))
}

func RenameProject(projectName string, kwargs map[string]interface{}, MAIN_DIR_PATH string) error {
	return files_manager.RenameDir(MAIN_DIR_PATH+"/"+projectName, kwargs["new_name"].(string))
}

func DeleteProject(projectName string, kwargs interface{}, MAIN_DIR_PATH string) error {
	return files_manager.DeleteDir(MAIN_DIR_PATH + "/" + projectName)
}
