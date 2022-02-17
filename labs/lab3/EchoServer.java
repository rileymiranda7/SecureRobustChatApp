import java.net.*;
import java.io.*;
import java.util.ArrayList;

public class EchoServer {
    static ThreadList threadlist = new ThreadList();

    public static void main(String[] args) throws IOException {
        if (args.length != 1) {
            System.err.println("Usage: java EchoServer <port number>");
            System.exit(1);
        }
        
        int portNumber = Integer.parseInt(args[0]);
        System.out.println("EchoServer by Riley Miranda");
        try {
            ServerSocket serverSocket =
                new ServerSocket(Integer.parseInt(args[0]));
            System.out.println("EchoServer is running at port " + Integer.parseInt(args[0]));
            while(true) {
                Socket clientSocket = serverSocket.accept(); 
                System.out.println("A client is connected ");    
                new EchoServerThread(threadlist, clientSocket).start();
            }
        } catch (IOException e) {
            System.out.println("Exception caught when trying to listen on port "
                + portNumber + " or listening for a connection");
            System.out.println(e.getMessage());
        }
    }
}
class EchoServerThread extends Thread{
    private Socket clientSocket = null;
    private ThreadList threadlist = null;
    public EchoServerThread(ThreadList threadlist, Socket socket) {
        clientSocket = socket;
        this.threadlist = threadlist;
    }
    public void run() {
        System.out.println("A new thread for client is running...");
        if(threadlist!=null) {
            threadlist.addThread(this);
            System.out.println("Inside thread: total clients: " + 
                threadlist.getNumberofThreads());
        }
        try {
            PrintWriter out =
                    new PrintWriter(clientSocket.getOutputStream(), true);                   
            BufferedReader in = new BufferedReader(
                    new InputStreamReader(clientSocket.getInputStream()));
        
            String inputLine;
            while ((inputLine = in.readLine()) != null) {
                System.out.println("received from client: " + inputLine);
                System.out.println("Echo back");
                out.println(inputLine);
            }
        } catch(IOException e) {
            System.out.println("IO Exception in thread" + e.getMessage());
        }
    }
}
class ThreadList{
    //private ... threadlist = //store the list of threads in this variable 
    private ArrayList<EchoServerThread> threadlist = new ArrayList<EchoServerThread>();
    public ThreadList(){        
    }
    public int getNumberofThreads(){
    //return the number of current threads
        return threadlist.size();
    }
    public void addThread(EchoServerThread newThread){
    //add the newthread object to the threadlist
    threadlist.add(newThread);    
    }
    public void removeThread(EchoServerThread thread){
    //remove the given thread from the threadlist 
    threadlist.remove(thread);      
    }
    public void sendToAll(String message){
    //ask each thread in the threadlist to send the given message to its client     
    }
}