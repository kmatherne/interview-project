<?php



	$reqBody = file_get_contents('php://input');
	$json = json_decode($reqBody, true);
	foreach ($json as $val) {
		print_r ($val);
	}



?>
