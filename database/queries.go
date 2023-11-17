package database

import (
	"cas-to-oauth2/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *Client) GenerateServiceTicket(st, service, username string, isDirect bool, expiration time.Time) {
	serviceTicketColl := c.Collection(constants.DB_COLLECTION_SERVICE_TICKETS)
	_, _ = serviceTicketColl.InsertOne(ctx, bson.M{"ticket": st, "service": service, "username": username, "isDirect": isDirect, "expires": expiration})
}

func (c *Client) GenerateTGT(tgt string, username string, expiration time.Time) {
	tgtColl := c.Collection(constants.DB_COLLECTION_TGT)
	_, _ = tgtColl.InsertOne(ctx, bson.M{"tgt": tgt, "username": username, "expires": expiration})
}

func (c *Client) ValidateServiceTicket(st, service string) (bool, string, bool) {
	serviceTicketColl := c.Collection(constants.DB_COLLECTION_SERVICE_TICKETS)
	filter := bson.M{"ticket": st, "service": service, "expires": bson.M{"$gte": time.Now()}}
	var result bson.M
	err := serviceTicketColl.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, "", false
		}
	}

	c.DeleteServiceTicket(st)

	return true, result["username"].(string), result["isDirect"].(bool)
}

func (c *Client) ValidateTGT(tgt string) (bool, string) {
	tgtColl := c.Collection(constants.DB_COLLECTION_TGT)
	filter := bson.M{"tgt": tgt, "expires": bson.M{"$gte": time.Now()}}
	var result bson.M
	err := tgtColl.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, ""
		}
	}

	return true, result["username"].(string)
}

func (c *Client) DeleteServiceTicket(st string) error {
	serviceTicketColl := c.Collection(constants.DB_COLLECTION_SERVICE_TICKETS)
	filter := bson.M{"ticket": st}
	_, err := serviceTicketColl.DeleteOne(ctx, filter)
	return err
}

func (c *Client) DeleteTGT(tgt string) error {
	tgtColl := c.Collection(constants.DB_COLLECTION_TGT)
	filter := bson.M{"tgt": tgt}
	_, err := tgtColl.DeleteOne(ctx, filter)
	return err
}
