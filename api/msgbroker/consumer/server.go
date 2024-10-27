package msgbroker

//type MsgBroker struct {
//	conf             *config.SystemConfig
//	db               *mongodb.DatabaseStorage
//	orderHandler     *order.OrderHandler
//	earnPointHandler *core_handle_point.CorePointHandler
//	log              *logger.Logger
//}
//
//func NewMsgBroker(
//	conf *config.SystemConfig,
//	orderHandler *order.OrderHandler,
//	earnPointHandler *core_handle_point.CorePointHandler,
//) *MsgBroker {
//	log := logger.GetLogger()
//	dbStorage, err := mongodb.ConnectMongoDB(context.Background(), &conf.MongoDBConfig)
//	if err != nil {
//		log.Fatal().Msgf("connect mongodb fail! %s", err)
//	}
//
//	return &MsgBroker{
//		conf:             conf,
//		db:               dbStorage,
//		orderHandler:     orderHandler,
//		earnPointHandler: earnPointHandler,
//		log:              log,
//	}
//}
//
//func (app *MsgBroker) Start(ctx context.Context) {
//	//Initialize the consumer for the receiver_create_order_success_dev topic
//	//csOrderSuccess := queuekafka.NewConsumer(app.conf.KafkaConfig, []string{app.conf.KafkaTopicConfig.TopicsRewardsPoint})
//	//csOrderSuccessEvent := order.NewConsumerOrder(csOrderSuccess, app.orderHandler)
//
//	//Initialize the consumer for the core_transaction_point_success_dev topic
//	csEarnPointSuccess := queuekafka.NewConsumer(app.conf.KafkaConfig, []string{app.conf.KafkaTopicConfig.TopicsCoreTransactionPointSuccess})
//	csEarnPointSuccessEvent := core_handle_point.NewConsumerEarnPoint(csEarnPointSuccess, app.earnPointHandler)
//
//	//go func() {
//	//	err := csOrderSuccessEvent.Start(ctx)
//	//	if err != nil {
//	//		app.log.Fatal().Err(err).Msg("Order success event consumer failed")
//	//	}
//	//}()
//
//	go func() {
//		err := csEarnPointSuccessEvent.Start(ctx)
//		if err != nil {
//			app.log.Fatal().Err(err).Msg("Earn point success event consumer failed")
//		}
//	}()
//}
