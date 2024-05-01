package playground;

import java.io.IOException;
import java.io.InputStream;
import java.net.HttpURLConnection;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.URL;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.Callable;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;

import org.junit.Test;

public class HTTPRequester {

	class Result {
		String id;
		String[] val;

		Result(String id, String[] val) {
			this.id = id;
			this.val = val;
		}
	}

	class MyCallable implements Callable<Result> {

		String id;

		MyCallable(String id) {
			this.id = id;
		}

		@Override
		public Result call() throws Exception {
			return new Result(id, req(id, id, id, id));
		}
	}

	void run(int numberOfRequests) {
		ExecutorService executor = Executors.newFixedThreadPool(10);

		List<Future<Result>> futureResults = new ArrayList<>();

		for (int i = 0; i < numberOfRequests; i++) {
			Future<Result> futureResult = executor.submit(new MyCallable(String.valueOf(i + 1)));
			futureResults.add(futureResult);
		}

		try {
			for (Future<Result> futureResult : futureResults) {
				Result result = futureResult.get();
				System.out.println("id " + result.id + " code " + result.val[0] + " content " + result.val[1]);
			}

		} catch (InterruptedException e) {
			Thread.currentThread().interrupt();
		} catch (ExecutionException e) {
			System.err.println("Task execution failed: " + e.getMessage());
		} finally {
			executor.shutdown();
		}
	}

	String[] req(String id, String lat, String lon, String query) {
		try {
			URL url = new URI("http://localhost:7000/app/memes?lat=" + lat + "&lon=" + lon + "&query=" + query).toURL();
			HttpURLConnection con = (HttpURLConnection) url.openConnection();
			con.setRequestMethod("GET");
			con.setRequestProperty("id", id);
			int responseCode = con.getResponseCode();
			InputStream inputStream = con.getInputStream();
			String content = new String(inputStream.readAllBytes(), StandardCharsets.UTF_8);
			return new String[] { String.valueOf(responseCode), content.trim() };
		} catch (IOException | URISyntaxException e) {
		}
		return new String[] { "", "" };
	}

	@Test
	public void test() {
		run(100);
	}

}
