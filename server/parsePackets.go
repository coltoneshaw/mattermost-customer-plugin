package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"

	"github.com/mattermost/mattermost/server/public/model"
	"gopkg.in/yaml.v3"
)

const (
	SupportPacketName = "support_packet.yaml"
	PluginFileName    = "plugins.json"
	ConfigFileName    = "sanitized_config.json"
)

func PostContainsSupportPackage(p *Plugin, post *model.Post) ([]*model.FileInfo, []string) {
	var supportPackets []*model.FileInfo
	var names []string

	for _, file := range post.FileIds {
		fileCheck, err := p.API.GetFileInfo(file)
		if err != nil {
			p.API.LogError("Failure checking for support packet." + err.Error())
		}
		supportPackets = append(supportPackets, fileCheck)
		names = append(names, fileCheck.Name)
	}

	return supportPackets, names
	// loop over each fileID
}

func unzipToMemory(zippedBytes []byte) ([]*model.FileData, error) {
	var fileContents []*model.FileData

	reader := bytes.NewReader(zippedBytes)
	zipReader, err := zip.NewReader(reader, int64(len(zippedBytes)))

	if err != nil {
		return nil, err
	}

	for _, file := range zipReader.File {
		// Open each file in the zip archive
		zippedFile, err := file.Open()
		if err != nil {
			return nil, err
		}

		defer zippedFile.Close()

		// Read the contents of the file into a byte slice
		contents, err := io.ReadAll(zippedFile)
		if err != nil {
			return nil, err
		}
		fmt.Println(filepath.Base(file.Name))
		// Store the file contents in the map using the file name as the key
		fileContents = append(fileContents, &model.FileData{
			Filename: filepath.Base(file.Name),
			Body:     contents,
		})
	}
	return fileContents, nil
}

func returnMarkdownResponse(packet *model.SupportPacket, config *model.Config, plugins *model.PluginsResponse) string {
	mdTable := "## Support Packet values\n"

	mdTable += "| Key | Value |\n| --- | --- |\n"
	mdTable += fmt.Sprintf("| %s | %v |\n", "Licensed To", packet.LicenseTo)
	mdTable += fmt.Sprintf("| %s | %v |\n", "Active Users", packet.ActiveUsers)
	mdTable += fmt.Sprintf("| %s | %v |\n", "Daily Active Users", packet.DailyActiveUsers)
	mdTable += fmt.Sprintf("| %s | %v |\n", "Monthly Active Users", packet.MonthlyActiveUsers)
	mdTable += fmt.Sprintf("| %s | %v |\n", "Server Arch", packet.ServerArchitecture)

	mdTable += fmt.Sprintf("| %s | %v |\n", "Server OS", packet.ServerOS)
	mdTable += fmt.Sprintf("| %s | %v |\n", "Server Version", packet.ServerVersion)
	mdTable += fmt.Sprintf("| %s | %v |\n", "Database Type", packet.DatabaseType)
	mdTable += fmt.Sprintf("| %s | %v |\n", "Database Version", packet.DatabaseVersion)

	mdTable += "\n"
	mdTable += "## Config Values\n"

	mdTable += "| Key | Value |\n| --- | --- |\n"
	mdTable += fmt.Sprintf("| %s | %t |\n", "High Availability", *config.ClusterSettings.Enable)
	mdTable += fmt.Sprintf("| %s | %t |\n", "SAML", *config.SamlSettings.Enable)
	mdTable += fmt.Sprintf("| %s | %t |\n", "LDAP", *config.LdapSettings.Enable)
	mdTable += fmt.Sprintf("| %s | %s |\n", "LDAP Groups", *config.LdapSettings.GroupFilter)
	mdTable += fmt.Sprintf("| %s | %t |\n", "Elasticsearch Search", *config.ElasticsearchSettings.EnableSearching)
	mdTable += fmt.Sprintf("| %s | %t |\n", "Elasticsearch Autocomplete", *config.ElasticsearchSettings.EnableAutocomplete)

	mdTable += "\n"
	mdTable += "## Plugin Info\n"

	mdTable += "| Plugin Name | Enabled | Version |\n| --- | :---: | --- |\n"

	for _, activePlugins := range plugins.Active {
		mdTable += fmt.Sprintf("| %s | %s | %s |\n", activePlugins.Name, ":white_check_mark:", activePlugins.Version)
	}
	for _, disabledPlugins := range plugins.Inactive {
		mdTable += fmt.Sprintf("| %s | %s | %s |\n", disabledPlugins.Name, "", disabledPlugins.Version)
	}

	fmt.Println(config)

	return mdTable
}

func unmarshalPacket(file *model.FileData) (*model.SupportPacket, error) {
	var packet *model.SupportPacket

	// Unmarshal the YAML into the struct
	err := yaml.Unmarshal(file.Body, &packet)
	if err != nil {
		return nil, err
	}

	return packet, nil
}

func unmarshalConfig(file *model.FileData) (*model.Config, error) {
	var config *model.Config

	// Unmarshal the YAML into the struct
	err := json.Unmarshal(file.Body, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func unmarshalPlugins(file *model.FileData) (*model.PluginsResponse, error) {
	var plugins *model.PluginsResponse

	// Unmarshal the YAML into the struct
	err := json.Unmarshal(file.Body, &plugins)
	if err != nil {
		return nil, err
	}

	return plugins, nil
}

// Responsible for downloading, reading, and processing the support packet
func ProcessSupportPackets(p *Plugin, packetArray []*model.FileInfo, post *model.Post) ([]*model.Post, error) {
	var posts []*model.Post

	// looking through all the packets in a post as you can upload more than one.
	for _, packet := range packetArray {
		fileData, err := p.API.GetFile(packet.Id)

		var packet *model.SupportPacket
		var config *model.Config
		var plugins *model.PluginsResponse

		if err != nil {
			return nil, err
		}

		unzippedFiles, err2 := unzipToMemory(fileData)
		if err2 != nil {
			p.API.LogError("Failure unpacking packet" + err2.Error())
			return nil, err2
		}

		// looking through everything that's in the zipped file to find
		for _, file := range unzippedFiles {
			var err error // Declare err outside the switch statement

			switch file.Filename {
			case SupportPacketName:
				data, e := unmarshalPacket(file)
				packet = data
				err = e
			case ConfigFileName:
				data, e := unmarshalConfig(file)
				config = data
				err = e
			case PluginFileName:
				data, e := unmarshalPlugins(file)
				err = e
				plugins = data
			}

			if err != nil {
				p.API.LogError("Error parsing support packet. Error:" + err.Error())
				return nil, err
			}
		}

		newPost, err := p.API.CreatePost(&model.Post{
			ChannelId: post.ChannelId,
			RootId:    post.Id,
			UserId:    p.botID,
			Message:   returnMarkdownResponse(packet, config, plugins),
		})
		if err != nil {
			p.API.LogError("Error parsing support packet. Error:" + err.Error())
			return nil, err
		}

		posts = append(posts, newPost)
	}

	return posts, nil
}
