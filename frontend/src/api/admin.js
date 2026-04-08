import { fetchJSON } from "./utils";
import { getApiPath } from "@/utils/url.js";

export function getStats() {
  return fetchJSON(getApiPath("api/admin/stats"));
}
