package data

import (
	"github.com/noptics/services/pkg/nproto"
)

// Store interface for the registry data
type Store interface {
	// SaveFiles should take the provided registry file and save it associated with the provided cluster/channel
	SaveFiles(cluster, channel string, files []*nproto.File) error
	// SetChannelMessage associates the root message for the channel
	SetChannelMessage(cluster, channel, message string) error
	// GetChannelData returns the root message and and proto files for the channel
	GetChannelData(cluster, channel string) (string, []*nproto.File, error)
	// SaveChannelData sets (overwrites) all the relevant data for a channel
	SaveChannelData(cluster, channel, message string, files []*nproto.File) error
	// GetChannels returns a list of channels the registry has data for
	GetChannels(cluster string) ([]string, error)
	// SaveCluster saves configuration data for a cluster
	SaveCluster(cluster *nproto.Cluster) (string, error)
	// GetCluster returns the configuration data for a cluster
	GetCluster(id string) (*nproto.Cluster, error)
	// GetClusters returns all cluster configurations
	GetClusters() ([]*nproto.Cluster, error)
}
