package app

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"gopkg.in/yaml.v2"
)

const (
	SupportPacketName = "support_packet.yaml"
	PluginFileName    = "plugins.json"
	ConfigFileName    = "sanitized_config.json"
)

func (s *customerService) MessageHasBeenPosted(post *model.Post) {
	if s.poster.IsFromPoster(post) || post.RootId != "" || len(post.FileIds) == 0 {
		return
	}

	supportPackets, names := postContainsSupportPackage(s, post)

	if len(supportPackets) == 0 {
		return
	}

	err := s.poster.PostMessageToThread(post.Id, &model.Post{
		ChannelId: post.ChannelId,
		Message:   "Uploading support packet for " + strings.Join(names, " ,"),
	})

	if err != nil {
		s.api.Log.Error("Failed in sending reply" + err.Error())
	}

	err = processSupportPackets(s, supportPackets, post)
	if err != nil {
		s.api.Log.Error("Failed processing packets" + err.Error())
	}
}

func postContainsSupportPackage(s *customerService, post *model.Post) ([]*model.FileInfo, []string) {
	var supportPackets []*model.FileInfo
	var names []string

	for _, id := range post.FileIds {
		fileCheck, err := s.api.File.GetInfo(id)
		if err != nil {
			s.api.Log.Error("Failure checking for support packet." + err.Error())
		}
		supportPackets = append(supportPackets, fileCheck)
		names = append(names, fileCheck.Name)
	}

	return supportPackets, names
}

func unzipToMemory(zippedBytes io.Reader) ([]*model.FileData, error) {
	var fileContents []*model.FileData

	// Read all data from the io.Reader into a byte slice
	data, err := io.ReadAll(zippedBytes)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(data)
	zipReader, err := zip.NewReader(reader, int64(len(data)))

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
	mdTable := "## Support Packet valuessss\n"

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
func processSupportPackets(s *customerService, packetArray []*model.FileInfo, post *model.Post) error {
	// looking through all the packets in a post as you can upload more than one.
	for _, packet := range packetArray {
		fileData, err := s.api.File.Get(packet.Id)

		var packet *model.SupportPacket
		var config *model.Config
		var plugins *model.PluginsResponse

		if err != nil {
			return err
		}

		unzippedFiles, err := unzipToMemory(fileData)
		if err != nil {
			s.api.Log.Error("Failure unpacking packet" + err.Error())
			return err
		}

		// looking through everything that's in the zipped file to find
		for _, file := range unzippedFiles {
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
				s.api.Log.Error("Error parsing support packet. Error:" + err.Error())
				return err
			}
		}

		// customerID, err := s.store.GetCustomerID(*config.ServiceSettings.SiteURL, packet.LicenseTo)

		// if err != nil {
		// 	s.api.Log.Error("Error getting customer ID. Error:" + err.Error())
		// 	return err
		// }

		err = s.poster.PostMessageToThread(post.Id, &model.Post{
			ChannelId: post.ChannelId,
			Message:   returnMarkdownResponse(packet, config, plugins),
		})

		if err != nil {
			s.api.Log.Error("Error parsing support packet. Error:" + err.Error())
			return err
		}
	}

	return nil
}
