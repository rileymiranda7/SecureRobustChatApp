<?php
ini_set('display_errors', 'On');
error_reporting(E_ALL);
	session_start();   
  $mysqli = new mysqli('localhost','mirandar1','1234','secad');
      if($mysqli->connect_errno){
        printf("Database connection failed: %s\n", $mysqli->connect_error);
        exit();
      } 
	if (securechecklogin($_POST["username"],$_POST["password"])) {
?>
	<h2> Welcome <?php echo htmlentities($_POST['username']); ?> !</h2>
<?php		
	}else{
		echo "<script>alert('Invalid username/password');</script>";
		die();
	}

  	function securechecklogin($username, $password) {
      global $mysqli;
  		$prepared_sql = "SELECT * FROM users WHERE username= ? " .
                      " AND password=password(?);";
  		if(!$stmt = $mysqli->prepare($prepared_sql))
        echo "Prepared statement error";
      $stmt->bind_param("ss", $username,$password);
      if(!$stmt->execute()) echo "Execute error";
      if(!$stmt->store_result()) echo "Store result error";
      $result = $stmt;
      if($result->num_rows ==1)
        return TRUE;
      return FALSE;
  	}
?>
