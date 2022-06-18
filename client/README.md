docker build --tag autogo-client-test .
docker run autogo-client-test

sudo docker inspect --format '{{ .NetworkSettings.IPAddress }}' 6cb599fe30ea


	//http.ListenAndServe(":8082", nil)
	//mosquitto_sub -h 7966f46d27eb47699c2b44744c43135c.s1.eu.hivemq.cloud -p 8883 -u autogo2 -P D3c3pt1c0n -t 'autogo/tank-01/move'
	//mosquitto_pub -h 7966f46d27eb47699c2b44744c43135c.s1.eu.hivemq.cloud -p 8883 -u autogo2 -P D3c3pt1c0n -t 'autogo/tank-01/move' -m 'Hello'