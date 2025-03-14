package regolith

import (
	"os"
	"path"

	"github.com/Bedrock-OSS/go-burrito/burrito"
)

type World struct {
	Id   string `json:"id"`   // The name of the world directory
	Name string `json:"name"` // The name of the world in levelname.txt
	Path string `json:"path"`
}

func ListWorlds(mojangDir string) ([]*World, error) {
	var worlds = make(map[string]World)
	var existingWorldNames = make(map[string]struct{}) // A set with duplicated world names
	var exists = struct{}{}

	worldsPath := path.Join(mojangDir, "minecraftWorlds")
	files, err := os.ReadDir(worldsPath)
	if err != nil {
		return nil, burrito.WrapErrorf(
			err, "Failed to list files in the directory.\nPath: %s", worldsPath)
	}
	for _, f := range files {
		if f.IsDir() {
			worldPath := path.Join(worldsPath, f.Name())
			worldname, err := os.ReadFile(
				path.Join(worldPath, "levelname.txt"))
			if err != nil {
				Logger.Warnf(
					"Unable to read levelname.txt from %q.", worldPath)
				continue
			}
			_, ok := existingWorldNames[string(worldname)]
			existingWorldNames[string(worldname)] = exists
			if ok { // The world with this name already exists
				delete(worlds, string(worldname))
				Logger.Warnf("Duplicated world name %q.", worldname)
				continue
			}
			worlds[string(worldname)] = World{
				Name: string(worldname),
				Id:   f.Name(),
				Path: worldPath,
			}
		}
	}
	// Convert result to list
	var result []*World
	for _, val := range worlds {
		result = append(result, &val)
	}
	return result, nil
}
