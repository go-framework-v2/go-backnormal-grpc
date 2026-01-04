package config

import "go.uber.org/zap"

type MongoDBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Db       string `mapstructure:"db"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func (c *MongoDBConfig) AppendFieldMap(fMap map[string]string) {
	fMap["mongodb.host"] = "MongoDB.Host"
	fMap["mongodb.port"] = "MongoDB.Port"
	fMap["mongodb.db"] = "MongoDB.Db"
	fMap["mongodb.username"] = "MongoDB.Username"
	fMap["mongodb.password"] = "MongoDB.Password"
}

func (c *MongoDBConfig) Print() {
	zap.L().Info("------------ mongodb ------------")
	zap.L().Info("-- ", zap.String("host", c.Host))
	zap.L().Info("-- ", zap.Int("port", c.Port))
	zap.L().Info("-- ", zap.String("db", c.Db))
	zap.L().Info("-- ", zap.String("username", c.Username))
	// zap.L().Info("-- ", zap.String("password", c.Password))
}
