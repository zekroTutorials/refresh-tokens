import EventEmitter from "events";
import { AccessToken, ErrorModel } from "./models";

const PREFIX =
  process.env.NODE_ENV === "development" ? "http://testserver/api" : "/api";

export default class RestAPI {
  public static readonly events = new EventEmitter();

  private static accessToken: string;

  public static login(username: string, password: string): Promise<any> {
    return this.post("auth/login", { username, password }, false);
  }

  public static me(): Promise<string> {
    return this.get("resources/me");
  }

  // ------------------------------------------------------------
  // --- HELPERS ---

  private static getAccessToken(): Promise<AccessToken> {
    return this.get<AccessToken>("auth/accesstoken");
  }

  private static get<T>(path: string, emitError: boolean = true): Promise<T> {
    return this.req<T>("GET", path, undefined, undefined, emitError);
  }

  private static post<T>(
    path: string,
    body?: any,
    emitError: boolean = true
  ): Promise<T> {
    return this.req<T>("POST", path, body, undefined, emitError);
  }

  private static delete<T>(
    path: string,
    emitError: boolean = true
  ): Promise<T> {
    return this.req<T>("DELETE", path, undefined, undefined, emitError);
  }

  private static async req<T>(
    method: string,
    path: string,
    body?: any,
    contentType: string | undefined = "application/json",
    emitError: boolean = true
  ): Promise<T> {
    let reqBody = undefined;
    if (body) {
      if (typeof body !== "string" && contentType === "application/json") {
        reqBody = JSON.stringify(body);
      } else {
        reqBody = body;
      }
    }

    const headers: { [key: string]: string } = {};
    if (contentType !== "multipart/form-data") {
      headers["content-type"] = contentType;
    }

    if (this.accessToken) {
      headers["authorization"] = "accessToken " + this.accessToken;
    }

    const res = await window.fetch(`${PREFIX}/${path}`, {
      method,
      headers,
      body: reqBody,
      credentials: "include",
    });

    if (res.status === 401) {
      try {
        const resBody = (await res.json()) as ErrorModel;
        if (resBody.error === "invalid access token") {
          const accessToken = await this.getAccessToken();
          this.accessToken = accessToken.token;
          return this.req(method, path, body, contentType, emitError);
        }
      } catch {}

      if (emitError) this.events.emit("authentication-error", res);
      throw new Error("auth error");
    }

    if (!res.ok) {
      if (emitError) this.events.emit("error", res);
      throw new Error(res.statusText);
    }

    if (res.status === 204 || res.headers.get("content-length") === "0") {
      return {} as T;
    }

    return res.json();
  }
}
