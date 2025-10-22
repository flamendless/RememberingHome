package main

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"

	"remembering-home/src/atlases"
	"remembering-home/src/common"
	"remembering-home/src/dialogues"
	"remembering-home/src/enums"
	"remembering-home/src/items"

	lua "github.com/yuin/gopher-lua"
)

func main() {
	if len(os.Args) < 3 {
		panic("Usage: lua_to_go --kind <atlas|item> --input <lua_file>")
	}

	var kind, inputFile string

	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--kind":
			if i+1 < len(os.Args) {
				kind = os.Args[i+1]
				i++
			} else {
				panic("Error: --kind requires a value (atlas or item)")
			}
		case "--input":
			if i+1 < len(os.Args) {
				inputFile = os.Args[i+1]
				i++
			} else {
				panic("Error: --input requires a filename")
			}
		}
	}

	if kind == "" {
		panic("Error: --kind is required")
	}
	if inputFile == "" {
		panic("Error: --input is required")
	}

	if kind != "atlas" && kind != "item" {
		panic("Error: --kind must be either 'atlas' or 'item'")
	}

	if !strings.HasSuffix(inputFile, ".lua") {
		panic("Error: input file must have .lua extension")
	}

	var outputFile string
	switch kind {
	case "atlas":
		dir := filepath.Base(filepath.Dir(inputFile))
		outputFile = filepath.Join("src/atlases", dir+".go")
	case "item":
		dir := filepath.Base(filepath.Dir(inputFile))
		outputFile = filepath.Join("src/items", dir+".go")
	}

	luaData, err := parseLuaFile(inputFile)
	if err != nil {
		panic(fmt.Sprintf("Error parsing Lua file: %v", err))
	}

	var goCode []byte
	var varName string
	switch kind {
	case "atlas":
		atlasData := convertToAtlasData(luaData)
		varName = generateVariableName("Atlas", inputFile)
		goCode, err = generateAtlasCode(atlasData, varName)
	case "item":
		itemData := convertToItemData(luaData)
		varName = generateVariableName("Item", inputFile)
		goCode, err = generateItemCode(itemData, varName)
	}
	if err != nil {
		panic(fmt.Sprintf("Error generating Go code: %v", err))
	}

	if err := os.WriteFile(outputFile, goCode, 0644); err != nil {
		panic(fmt.Sprintf("Error writing output file: %v", err))
	}

	fmt.Printf("Successfully converted %s to %s\n", inputFile, outputFile)
}

func parseLuaFile(filename string) (*lua.LTable, error) {
	L := lua.NewState()
	defer L.Close()
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if err := L.DoString(string(content)); err != nil {
		return nil, err
	}
	ret := L.Get(-1)
	if ret.Type() != lua.LTTable {
		return nil, fmt.Errorf("expected table, got %s", ret.Type())
	}
	return ret.(*lua.LTable), nil
}

func convertToAtlasData(luaTable *lua.LTable) atlases.AtlasData {
	var frames []atlases.FrameData
	var metadata atlases.Metadata

	luaTable.ForEach(func(key, value lua.LValue) {
		switch key.String() {
		case "frames":
			if value.Type() == lua.LTTable {
				framesTable := value.(*lua.LTable)
				framesTable.ForEach(func(frameKey, frameValue lua.LValue) {
					if frameValue.Type() == lua.LTTable {
						frameTable := frameValue.(*lua.LTable)
						frame := convertFrameTable(frameKey.String(), frameTable)
						frames = append(frames, frame)
					}
				})
			}
		case "meta":
			if value.Type() == lua.LTTable {
				metadata = convertMetaTable(value.(*lua.LTable))
			}
		}
	})

	return atlases.AtlasData{
		Frames:   frames,
		Metadata: metadata,
	}
}

func convertFrameTable(id string, frameTable *lua.LTable) atlases.FrameData {
	if id == "" {
		panic("Error: empty item ID")
	}

	itemID := enums.Item(id)
	if itemID.Constant() == "ItemUndefined" {
		panic(fmt.Sprintf("Error: undefined item ID '%s'", id))
	}

	frame := atlases.FrameData{
		ID:   itemID,
		Pos:  common.Vec2{X: 0, Y: 0},
		Size: common.Vec2{X: 0, Y: 0},
	}

	frameTable.ForEach(func(key, value lua.LValue) {
		switch key.String() {
		case "x":
			if value.Type() == lua.LTNumber {
				frame.Pos.X = float64(value.(lua.LNumber))
			}
		case "y":
			if value.Type() == lua.LTNumber {
				frame.Pos.Y = float64(value.(lua.LNumber))
			}
		case "w":
			if value.Type() == lua.LTNumber {
				frame.Size.X = float64(value.(lua.LNumber))
			}
		case "h":
			if value.Type() == lua.LTNumber {
				frame.Size.Y = float64(value.(lua.LNumber))
			}
		}
	})

	return frame
}

func convertMetaTable(metaTable *lua.LTable) atlases.Metadata {
	metadata := atlases.Metadata{
		Padding:   0,
		Extrude:   0,
		Size:      common.Vec2{X: 0, Y: 0},
		QuadCount: 0,
	}

	metaTable.ForEach(func(key, value lua.LValue) {
		switch key.String() {
		case "padding":
			if value.Type() == lua.LTNumber {
				metadata.Padding = int(value.(lua.LNumber))
			}
		case "extrude":
			if value.Type() == lua.LTNumber {
				metadata.Extrude = int(value.(lua.LNumber))
			}
		case "atlasWidth":
			if value.Type() == lua.LTNumber {
				metadata.Size.X = float64(value.(lua.LNumber))
			}
		case "atlasHeight":
			if value.Type() == lua.LTNumber {
				metadata.Size.Y = float64(value.(lua.LNumber))
			}
		case "quadCount":
			if value.Type() == lua.LTNumber {
				metadata.QuadCount = int(value.(lua.LNumber))
			}
		}
	})

	return metadata
}

func convertToItemData(luaTable *lua.LTable) []items.ItemData {
	var goData []items.ItemData

	luaTable.ForEach(func(key, value lua.LValue) {
		if value.Type() != lua.LTTable {
			return
		}

		itemTable := value.(*lua.LTable)
		data := convertTableItem(itemTable)
		goData = append(goData, data)
	})

	return goData
}

func convertTableItem(itemTable *lua.LTable) items.ItemData {
	data := items.ItemData{
		Item:      enums.ItemUndefined,
		Name:      "",
		Pos:       common.Vec2{X: 0, Y: 0},
		NoCol:     false,
		ReqColDir: enums.DirectionUndefined,
		Tags:      []enums.ItemClass{},
		Dialogue:  dialogues.DialogueKeys{},
	}

	itemTable.ForEach(func(key, value lua.LValue) {
		switch key.String() {
		case "id":
			if value.Type() == lua.LTString {
				id := value.String()
				if id == "" {
					panic("Error: empty item ID")
				}
				itemID := enums.Item(id)
				if itemID.Constant() == "ItemUndefined" {
					panic(fmt.Sprintf("Error: undefined item ID '%s'", id))
				}
				data.Item = itemID
			}
		case "name":
			if value.Type() == lua.LTString {
				data.Name = value.String()
			}
		case "x":
			if value.Type() == lua.LTNumber {
				data.Pos.X = float64(value.(lua.LNumber))
			}
		case "y":
			if value.Type() == lua.LTNumber {
				data.Pos.Y = float64(value.(lua.LNumber))
			}
		case "no_col":
			if value.Type() == lua.LTBool {
				data.NoCol = bool(value.(lua.LBool))
			}
		case "req_col_dir":
			if value.Type() == lua.LTNumber {
				dir := int(value.(lua.LNumber))
				switch dir {
				case -1:
					data.ReqColDir = enums.DirectionLeft
				case 1:
					data.ReqColDir = enums.DirectionRight
				default:
					data.ReqColDir = enums.DirectionUndefined
				}
			}
		case "is_door":
			if value.Type() == lua.LTBool && bool(value.(lua.LBool)) {
				data.Tags = append(data.Tags, enums.ItemClassDoor)
			}
		case "dialogue":
			if value.Type() == lua.LTTable {
				dialogueTable := value.(*lua.LTable)
				if dialogueTable.Len() >= 2 {
					roomID := dialogueTable.RawGetInt(1)
					itemID := dialogueTable.RawGetInt(2)
					if roomID.Type() == lua.LTString && itemID.Type() == lua.LTString {
						roomStr := roomID.String()
						itemStr := itemID.String()

						if roomStr == "" {
							panic("Error: empty room ID")
						}
						if itemStr == "" {
							panic("Error: empty item ID")
						}

						room := enums.Room(roomStr)
						item := enums.Item(itemStr)

						if room.Constant() == "RoomUndefined" {
							panic(fmt.Sprintf("Error: undefined room ID '%s'", roomStr))
						}
						if item.Constant() == "ItemUndefined" {
							panic(fmt.Sprintf("Error: undefined item ID '%s'", itemStr))
						}

						data.Dialogue = dialogues.DialogueKeys{
							Room: room,
							Item: item,
						}
					}
				}
			}
		}
	})

	return data
}

func generateAtlasCode(data atlases.AtlasData, varName string) ([]byte, error) {
	var buf strings.Builder

	buf.WriteString("// Code generated by lua_to_go. DO NOT EDIT.\n")
	buf.WriteString("package atlases\n\n")

	buf.WriteString("import (\n")
	buf.WriteString("\t\"remembering-home/src/common\"\n")
	buf.WriteString("\t\"remembering-home/src/enums\"\n")
	buf.WriteString(")\n\n")

	buf.WriteString(fmt.Sprintf("var %s = AtlasData{\n", varName))

	buf.WriteString("\tFrames: []FrameData{\n")
	for _, frame := range data.Frames {
		buf.WriteString("\t\t{\n")
		buf.WriteString(fmt.Sprintf("\t\t\tID: enums.%s,\n", frame.ID.Constant()))
		buf.WriteString(fmt.Sprintf("\t\t\tPos: common.Vec2{X: %g, Y: %g},\n", frame.Pos.X, frame.Pos.Y))
		buf.WriteString(fmt.Sprintf("\t\t\tSize: common.Vec2{X: %g, Y: %g},\n", frame.Size.X, frame.Size.Y))
		buf.WriteString("\t\t},\n")
	}
	buf.WriteString("\t},\n")

	buf.WriteString("\tMetadata: Metadata{\n")
	buf.WriteString(fmt.Sprintf("\t\tPadding: %d,\n", data.Metadata.Padding))
	buf.WriteString(fmt.Sprintf("\t\tExtrude: %d,\n", data.Metadata.Extrude))
	buf.WriteString(fmt.Sprintf("\t\tSize: common.Vec2{X: %g, Y: %g},\n", data.Metadata.Size.X, data.Metadata.Size.Y))
	buf.WriteString(fmt.Sprintf("\t\tQuadCount: %d,\n", data.Metadata.QuadCount))
	buf.WriteString("\t},\n")

	buf.WriteString("}\n")

	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		return nil, err
	}

	return formatted, nil
}

func generateItemCode(data []items.ItemData, varName string) ([]byte, error) {
	var buf strings.Builder

	buf.WriteString("// Code generated by lua_to_go. DO NOT EDIT.\n")
	buf.WriteString("package items\n\n")

	buf.WriteString("import (\n")
	buf.WriteString("\t\"remembering-home/src/common\"\n")
	buf.WriteString("\t\"remembering-home/src/dialogues\"\n")
	buf.WriteString("\t\"remembering-home/src/enums\"\n")
	buf.WriteString(")\n\n")

	buf.WriteString(fmt.Sprintf("var %s = []ItemData{\n", varName))

	for _, item := range data {
		buf.WriteString("\t{\n")
		buf.WriteString(fmt.Sprintf("\t\tItem: enums.%s,\n", item.Item.Constant()))
		if item.Name != "" {
			buf.WriteString(fmt.Sprintf("\t\tName: \"%s\",\n", item.Name))
		}
		buf.WriteString(fmt.Sprintf("\t\tPos: common.Vec2{X: %g, Y: %g},\n", item.Pos.X, item.Pos.Y))
		if item.NoCol {
			buf.WriteString("\t\tNoCol: true,\n")
		}
		if item.ReqColDir != enums.DirectionUndefined {
			var dirStr string
			switch item.ReqColDir {
			case enums.DirectionLeft:
				dirStr = "DirectionLeft"
			case enums.DirectionRight:
				dirStr = "DirectionRight"
			default:
				dirStr = "DirectionUndefined"
			}
			buf.WriteString(fmt.Sprintf("\t\tReqColDir: enums.%s,\n", dirStr))
		}
		if len(item.Tags) > 0 {
			buf.WriteString("\t\tTags: []enums.ItemClass{")
			for i, tag := range item.Tags {
				if i > 0 {
					buf.WriteString(", ")
				}
				buf.WriteString(fmt.Sprintf("enums.%s", tag.Constant()))
			}
			buf.WriteString("},\n")
		}
		if item.Dialogue.Room != enums.RoomUndefined || item.Dialogue.Item != enums.ItemUndefined {
			buf.WriteString("\t\tDialogue: dialogues.DialogueKeys{\n")
			if item.Dialogue.Room != enums.RoomUndefined {
				buf.WriteString(fmt.Sprintf("\t\t\tRoom: enums.%s,\n", item.Dialogue.Room.Constant()))
			}
			if item.Dialogue.Item != enums.ItemUndefined {
				buf.WriteString(fmt.Sprintf("\t\t\tItem: enums.%s,\n", item.Dialogue.Item.Constant()))
			}
			buf.WriteString("\t\t},\n")
		}
		buf.WriteString("\t},\n")
	}

	buf.WriteString("}\n")

	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		return nil, err
	}

	return formatted, nil
}

func generateVariableName(kind string, inputFile string) string {
	dir := filepath.Base(filepath.Dir(inputFile))
	parts := strings.Split(dir, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return kind + strings.Join(parts, "")
}
