# delivery-service
The delivery service sends notifications to the users. It listens to the messages in the redis channel that are
sent my the fanout-consumer and sends it to the user if they are connected via web-socket.
