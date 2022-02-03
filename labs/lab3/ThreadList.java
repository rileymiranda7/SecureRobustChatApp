class ThreadList{
	//private ... threadlist = //store the list of threads in this variable	
	public ThreadList(){		
	}
	public int getNumberofThreads(){
	//return the number of current threads
	}
	public void addThread(EchoServerThread newthread){
	//add the newthread object to the threadlist	
	}
	public void removeThread(EchoServerThread thread){
	//remove the given thread from the threadlist		
	}
	public void sendToAll(String message){
	//ask each thread in the threadlist to send the given message to its client		
	}
}