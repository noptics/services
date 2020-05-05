package dynamo

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
	"github.com/noptics/services/pkg/golog"
	"github.com/noptics/services/pkg/nproto"
)

const (
	CHANNELS_TABLE             = "noptics_registry_channels"
	CHANNELS_FIELD_CLUSTER     = "clusterID"
	CHANNELS_FIELD_CHANNEL     = "channel"
	CHANNELS_FIELD_FILES       = "files"
	CHANNELS_FIELD_MESSAGE     = "message"
	CLUSTERS_TABLE             = "noptics_registry_clusters"
	CLUSTERS_FIELD_ID          = "clusterID"
	CLUSTERS_FIELD_NAME        = "clusterName"
	CLUSTERS_FIELD_DESCRIPTION = "clusterDescription"
	CLUSTERS_FIELD_SERVER      = "servers"
	CLUSTERS_FIELD_DATA        = "data"
)

// dynamo implements the store interface for dynamodb
type Dynamo struct {
	// awsSession holds access credentials to the db
	awsSession *session.Session
	// endpoint is used for testing. Empty is the default and will connect to live dynamodb servers
	endpoint string
	// tablePrefix is used to segregate data by env
	tablePrefix string
}

func New(config map[string]string) (*Dynamo, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("unable to create aws session (%s)", err.Error())
	}

	return &Dynamo{
		awsSession:  sess,
		endpoint:    config["endpoint"],
		tablePrefix: config["prefix"],
	}, nil
}

func (db *Dynamo) Table(name string) *string {
	return aws.String(db.tablePrefix + name)
}

func dbClient(db *Dynamo) *dynamodb.DynamoDB {
	var ep *string
	if len(db.endpoint) == 0 {
		ep = nil
	} else {
		ep = aws.String(db.endpoint)
	}
	return dynamodb.New(db.awsSession, &aws.Config{Endpoint: ep})
}

func (db *Dynamo) SaveFiles(cluster, channel string, files []*nproto.File) error {
	update := expression.Set(
		expression.Name(CHANNELS_FIELD_FILES),
		expression.Value(files),
	)

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return err
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Key: map[string]*dynamodb.AttributeValue{
			CHANNELS_FIELD_CLUSTER: {
				S: aws.String(cluster),
			},
			CHANNELS_FIELD_CHANNEL: {
				S: aws.String(channel),
			},
		},
		TableName:        db.Table(CHANNELS_TABLE),
		UpdateExpression: expr.Update(),
	}

	golog.Debugw("save files data", "files", files, "channel", channel, "cluster", cluster, "input", input)

	_, err = dbClient(db).UpdateItem(input)

	return err
}

func (db *Dynamo) GetChannelData(cluster, channel string) (string, []*nproto.File, error) {
	client := dbClient(db)

	golog.Debugw("get files", "cluster", cluster, "channel", channel)

	resp, err := client.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			CHANNELS_FIELD_CLUSTER: {
				S: aws.String(cluster),
			},
			CHANNELS_FIELD_CHANNEL: {
				S: aws.String(channel),
			},
		},
		TableName: db.Table(CHANNELS_TABLE),
	})

	if err != nil {
		return "", nil, err
	}

	if len(resp.Item) == 0 {
		return "", nil, nil
	}

	f := []*nproto.File{}

	err = dynamodbattribute.Unmarshal(resp.Item[CHANNELS_FIELD_FILES], &f)
	if err != nil {
		return "", nil, err
	}

	msg := ""
	if att, ok := resp.Item[CHANNELS_FIELD_MESSAGE]; ok {
		msg = *att.S
	}

	return msg, f, nil
}

func (db *Dynamo) SetChannelMessage(cluster, channel, message string) error {
	update := expression.Set(
		expression.Name(CHANNELS_FIELD_MESSAGE),
		expression.Value(message),
	)

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return err
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Key: map[string]*dynamodb.AttributeValue{
			CHANNELS_FIELD_CLUSTER: {
				S: aws.String(cluster),
			},
			CHANNELS_FIELD_CHANNEL: {
				S: aws.String(channel),
			},
		},
		TableName:        db.Table(CHANNELS_TABLE),
		UpdateExpression: expr.Update(),
	}

	golog.Debugw("set channel message", "channel", channel, "cluster", cluster, "message", message, "input", input)

	_, err = dbClient(db).UpdateItem(input)

	return err
}

func (db *Dynamo) GetChannels(cluster string) ([]string, error) {
	keyExp := expression.Key(CHANNELS_FIELD_CLUSTER).Equal(expression.Value(cluster))
	proj := expression.NamesList(expression.Name(CHANNELS_FIELD_CHANNEL))

	exp, err := expression.NewBuilder().WithKeyCondition(keyExp).WithProjection(proj).Build()
	if err != nil {
		return nil, err
	}

	client := dbClient(db)
	resp, err := client.Query(&dynamodb.QueryInput{
		ExpressionAttributeValues: exp.Values(),
		ExpressionAttributeNames:  exp.Names(),
		KeyConditionExpression:    exp.KeyCondition(),
		ProjectionExpression:      exp.Projection(),
		TableName:                 db.Table(CHANNELS_TABLE),
	})

	if err != nil {
		return nil, err
	}

	chans := []string{}

	golog.Debugw("query items", "cluster", cluster, "items", resp.Items)

	for _, i := range resp.Items {
		chans = append(chans, *i["channel"].S)
	}

	return chans, nil
}

func (db *Dynamo) SaveChannelData(cluster, channel, message string, files []*nproto.File) error {
	update := expression.Set(
		expression.Name(CHANNELS_FIELD_MESSAGE),
		expression.Value(message),
	).Set(
		expression.Name(CHANNELS_FIELD_FILES),
		expression.Value(files),
	)

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return err
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Key: map[string]*dynamodb.AttributeValue{
			CHANNELS_FIELD_CLUSTER: {
				S: aws.String(cluster),
			},
			CHANNELS_FIELD_CHANNEL: {
				S: aws.String(channel),
			},
		},
		TableName:        db.Table(CHANNELS_TABLE),
		UpdateExpression: expr.Update(),
	}

	golog.Debugw("set channel data", "channel", channel, "cluster", cluster, "message", message, "files", files, "input", input)

	_, err = dbClient(db).UpdateItem(input)

	return err
}

func (db *Dynamo) SaveCluster(cluster *nproto.Cluster) (string, error) {
	if cluster.Id == "" {
		cluster.Id = uuid.New().String()
	}

	data, err := dynamodbattribute.Marshal(cluster)
	if err != nil {
		return "", err
	}

	update := expression.Set(
		expression.Name(CLUSTERS_FIELD_ID),
		expression.Value(cluster.Id),
	).Set(
		expression.Name(CLUSTERS_FIELD_DATA),
		expression.Value(data),
	)

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return "", err
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Key: map[string]*dynamodb.AttributeValue{
			CLUSTERS_FIELD_ID: {
				S: aws.String(cluster.Id),
			},
		},
		TableName:        db.Table(CLUSTERS_TABLE),
		UpdateExpression: expr.Update(),
	}

	golog.Debugw("save cluster data", "clustser", cluster, "input", input)

	_, err = dbClient(db).UpdateItem(input)

	return cluster.Id, err
}

func (db *Dynamo) GetCluster(id string) (*nproto.Cluster, error) {
	item, err := dbClient(db).GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			CLUSTERS_FIELD_ID: {
				S: aws.String(id),
			},
		},
		TableName: db.Table(CLUSTERS_TABLE),
	})

	if err != nil {
		return nil, err
	}

	cluster := &nproto.Cluster{}

	if data, ok := item.Item[CLUSTERS_FIELD_DATA]; !ok {
		return nil, fmt.Errorf("could not find cluster %s", id)
	} else {
		err = dynamodbattribute.Unmarshal(data, cluster)
	}

	return cluster, err
}

func (db *Dynamo) GetClusters() ([]*nproto.Cluster, error) {
	proj := expression.NamesList(expression.Name(CLUSTERS_FIELD_DATA))

	exp, err := expression.NewBuilder().WithProjection(proj).Build()
	if err != nil {
		return nil, err
	}

	resp, err := dbClient(db).Scan(&dynamodb.ScanInput{
		TableName:            db.Table(CLUSTERS_TABLE),
		ProjectionExpression: exp.Projection(),
	})

	if err != nil {
		return nil, err
	}

	clusters := []*nproto.Cluster{}

	for _, i := range resp.Items {
		c := &nproto.Cluster{}
		err = dynamodbattribute.Unmarshal(i[CLUSTERS_FIELD_DATA], c)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling cluster")
		}
		clusters = append(clusters, c)
	}

	return clusters, nil
}
