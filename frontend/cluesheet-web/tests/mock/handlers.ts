import { http, HttpResponse } from "msw";
import { mockCluesheetUUID } from "../lib/data";

export default [
  http.get(`/api/v1/cluesheet/${mockCluesheetUUID}`, async ({ request }) => {
    console.debug("Hello from mocked join API.");
    // OK we're chilling. Return 200
    return HttpResponse.json({ detail: "CSH" }, { status: 201 });
  }),
  /*
  http.post("/api/v1/nn-assign/", async ({ request }) => {
    console.debug("Hello from mocked NN Assign API.");

    const requestJson = await request.json();

    return HttpResponse.json(
      { detail: "Hello"},
      { status: 200 },
    );
  }),
  */
];
