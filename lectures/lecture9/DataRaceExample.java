public class DataRaceExample extends Thread {
	private static int shared = 0;
	int id;
	DataRaceExample(int id){
		this.id = id;
	}
	public void run(){
		System.out.println("In Thread " + id + ", shared= " + shared);
		shared++;
	}
	public static void main(String args[]){
		for (int i =0 ; i < 5; i ++){
			Thread thread = new DataRaceExample(i);
			thread.start();
		}
	}
}
