import { AxiosError } from "axios";

export const parseApiError = (err: any): HttpError => {
  const error = err as AxiosError
  if (error.response) {
    // The request was made and the server responded with a non-2xx status code
    return new HttpError(err.response.data.status, err.response.data.message, err.response.data.error);
  } else if (err.request) {
    // The request was made but no response was received
    return new HttpError(500, "Network error or no response from server", null);
  } else {
    // Something else happened while setting up the request
    return new HttpError(500, `Unexpected error: ${err.message}`, null);
  }
}

export class HttpError extends Error {
  public status: number;
  public error: any;

  constructor(status: number, message: string, error: any) {
    super(message)
    this.name = 'HttpError';
    this.status = status;
    this.error = error
  }
}