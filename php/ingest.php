 <?php



	$reqBody = file_get_contents('php://input');
	$json = json_decode($reqBody, true);


	$redis = new Redis();
	$connected = $redis->connect('127.0.0.1', 8888);
	$auth = $redis->auth("siege87751");
	if ($connected) {
		$id = $redis->get('next_endpoint_id');
		print $id;
		if (!$id) {
			print "\nSetting Id\n";
			$id = 1000;
			$redis->set('next_endpoint_id', $id);
		}
		if (isset($json['data']) && isset($json['endpoint'])) {
			$redis->hMSet('endpoint:' . $id, $json['endpoint']);
			foreach ($json as $val) {
				print_r ($val);
			}
			$redis->hMSet('data:' . $id, $json['data']);
			$pushed = $redis->rPush('endpoint_query', $id);
			if (!$pushed) {
				print "\nError pushing to enpoint query\n";
			}
			$redis->incr('next_endpoint_id');
		}
	} else {
		print "\nCouldn't connect to Redis.\n";
	}



?>
