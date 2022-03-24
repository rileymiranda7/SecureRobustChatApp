<?php
	session_start();    
	if (checklogin_sql($_POST["username"],$_POST["password"])) {
?>
	<h2> Welcome <?php echo $_POST['username']; ?> !</h2>
<?php		
	}else{
		echo "<script>alert('Invalid username/password');</script>";
		die();
	}
	/*function checklogin($username, $password) {
		$account = array("admin","1234");
		if (($username== $account[0]) and ($password == $account[1])) 
		  return TRUE;
		else return FALSE;
  	}*/
  	function checklogin_sql($username, $password) {
  		$mysqli = new my sqli('locahlhost',
  								'mirandar1', //Database username
  								'Peter6696', //Database password
  								'secad'); // Name of database
  		if($mysqli->connect_errno){
  			printf("Database connection failed: %s\n", $mysqli->connect_error);
  			exit();
  		}
  		$sql = "SELECT * FROM users WHERE username='" . $username. "' ";
  		$sql = $sql . "AND password = password('" . $password . "')";
  		echo "DEBUG>sql= $sql";
  		return TRUE;
		/*$account = array("admin","1234");
		if (($username== $account[0]) and ($password == $account[1])) 
		  return TRUE;
		else return FALSE;*/
  	}
?>
