package converter

import (
	"fmt"

	"ajpa/analyzer"

	"github.com/azure/azure-dev/cli/azd/pkg/project"
)

func ProjectAnalysisResultToAzdProjectConfig(result analyzer.ProjectAnalysisResult) (project.ProjectConfig, error) {
	config := project.ProjectConfig{
		Name: result.Name,
	}
	config.Services = make(map[string]*project.ServiceConfig)
	for name, app := range result.Applications {
		config.Services[name] = &project.ServiceConfig{
			Project:      &config,
			Name:         name,
			RelativePath: app.ProjectRelativePath,
			Host:         project.ContainerAppTarget, // todo: support other kinds.
			Language:     project.ServiceLanguageJava,
		}
	}
	config.Resources = make(map[string]*project.ResourceConfig)
	for name, service := range result.Services {
		resourceType, err := toResourceType(service)
		if err != nil {
			return project.ProjectConfig{}, err
		}
		props, err := toProps(service)
		if err != nil {
			return project.ProjectConfig{}, err
		}
		config.Resources[name] = &project.ResourceConfig{
			Project: &config,
			Name:    name,
			Type:    resourceType,
			Props:   props,
		}
	}
	for appName, serviceNameMap := range result.ApplicationToBackingService {
		hostingName := result.ApplicationToHostingService[appName]
		for serviceName := range serviceNameMap {
			config.Resources[hostingName].Uses = append(config.Resources[hostingName].Uses, serviceName)
		}
	}
	return config, nil
}

func toResourceType(service analyzer.Service) (project.ResourceType, error) {
	switch service.(type) {
	case analyzer.AzureContainerApp:
		return project.ResourceTypeHostContainerApp, nil
	case analyzer.AzureDatabaseForPostgresql:
		return project.ResourceTypeDbPostgres, nil
	case analyzer.AzureDatabaseForMysql:
		return project.ResourceTypeDbPostgres, nil // todo: change to mysql when azd support mysql
	case analyzer.AzureCacheForRedis:
		return project.ResourceTypeDbRedis, nil
	case analyzer.AzureCosmosDbForMongoDb:
		return project.ResourceTypeDbMongo, nil
	case analyzer.AzureCosmosDb:
		return project.ResourceTypeDbMongo, nil // todo: change to cosmos when azd support cosmos
	case analyzer.AzureServiceBus:
		return project.ResourceTypeMessagingServiceBus, nil
	case analyzer.AzureEventHubs:
		return project.ResourceTypeMessagingEventHubs, nil
	case analyzer.AzureStorageAccount:
		return project.ResourceTypeStorage, nil
	default:
		return "", fmt.Errorf("unknown service type: %v", service)
	}
}

func toProps(service analyzer.Service) (interface{}, error) {
	switch s := service.(type) {
	case analyzer.AzureContainerApp:
		return project.ContainerAppProps{
			Port: 8080, // todo: support non-web app.
		}, nil
	case analyzer.AzureDatabaseForPostgresql, // todo: Add database name in PostgresqlProps
		analyzer.AzureDatabaseForMysql,
		analyzer.AzureCacheForRedis,
		analyzer.AzureCosmosDbForMongoDb,
		analyzer.AzureCosmosDb:
		return nil, nil
	case analyzer.AzureServiceBus:
		return project.ServiceBusProps{
			Queues: s.Queues,
			Topics: s.Topics,
		}, nil
	case analyzer.AzureEventHubs:
		return project.EventHubsProps{
			Hubs: s.Hubs,
		}, nil
	case analyzer.AzureStorageAccount:
		return project.StorageProps{
			Containers: s.Containers,
		}, nil
	default:
		return "", fmt.Errorf("unknown service type when get Props: %v", service)
	}
}
