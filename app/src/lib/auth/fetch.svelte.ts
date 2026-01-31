import { API_BASE_URL } from '$lib/constants';
import { FirebaseAuthentication } from '@capacitor-firebase/authentication';
import type { paths, ValidationError } from '$lib/types/api';

async function fetchWithAuth(url: string, options: RequestInit = {}) {
	const result = await FirebaseAuthentication.getIdToken();
	const token = result.token;

	const headers = {
		...options.headers,
		Authorization: token ? `Bearer ${token}` : '',
		'Content-Type': 'application/json'
	};

	const response = await fetch(API_BASE_URL + url, {
		...options,
		headers
	});

	return response;
}

// ---------- Typed helpers built from OpenAPI ----------

type HttpMethod = 'get' | 'post' | 'put' | 'delete' | 'patch';

export type GetResponseData<
	P extends keyof paths,
	M extends HttpMethod = 'get'
> = paths[P][M] extends { responses: { 200: { content: { 'application/json': infer T } } } }
	? T
	: paths[P][M] extends { responses: { 201: { content: { 'application/json': infer T } } } }
		? T
		: paths[P][M] extends { responses: { 200: { content?: never } } }
			? never
			: never;

export type GetRequestData<
	P extends keyof paths,
	M extends HttpMethod = 'post'
> = paths[P][M] extends { requestBody: { content: { '*/*': infer T } } }
	? T
	: paths[P][M] extends { requestBody: { content: { 'application/json': infer T } } }
		? T
		: never;

export type GetPathParams<P extends keyof paths> = paths[P] extends {
	parameters: { path: infer T };
}
	? T extends Record<string, unknown>
		? { [K in keyof T]: T[K] }
		: never
	: never;

type FetchSuccess<P extends keyof paths, M extends HttpMethod = 'get'> = {
	ok: true;
	loading: false;
	data: GetResponseData<P, M>;
};

type FetchFailure = {
	ok: false;
	loading: false;
	error: {
		status?: number;
		message?: string;
		errors?: ValidationError[];
		body?: unknown;
	};
};

export type FetchResult<P extends keyof paths, M extends HttpMethod = 'get'> =
	| FetchSuccess<P, M>
	| FetchFailure;

function applyPathParams<P extends keyof paths>(apiPath: P, params: GetPathParams<P> | undefined) {
	let path = apiPath as string;
	if (params && typeof params === 'object') {
		for (const [key, value] of Object.entries(params)) {
			path = path.replace(`{${key}}`, encodeURIComponent(String(value ?? '')));
		}
	}
	return path;
}

async function parseJsonSafe(res: Response): Promise<unknown | undefined> {
	const text = await res.text();
	if (!text) return undefined;
	try {
		return JSON.parse(text);
	} catch {
		return undefined;
	}
}

export async function fetchData<P extends keyof paths, M extends HttpMethod = 'get'>(
	apiPath: P,
	options: GetPathParams<P> extends never
		? {
				method?: M;
				requestOptions?: RequestInit;
				requestName?: string;
				noError?: boolean;
				body?: GetRequestData<P, M>;
				params?: never;
			}
		: {
				method?: M;
				requestOptions?: RequestInit;
				requestName?: string;
				noError?: boolean;
				body?: GetRequestData<P, M>;
				params: GetPathParams<P>;
			}
): Promise<FetchResult<P, M>> {
	const method = (options?.method ?? ('get' as M)) as HttpMethod;
	const params = (options as { params?: GetPathParams<P> }).params;
	const path = applyPathParams(apiPath, params);

	try {
		const reqInit: RequestInit = {
			method: method.toUpperCase(),
			...options.requestOptions
		};

		if (options.body !== undefined) {
			(reqInit.headers as Record<string, string>) = {
				...(reqInit.headers as Record<string, string>),
				'Content-Type': 'application/json'
			};
			reqInit.body = JSON.stringify(options.body);
		}

		const res = await fetchWithAuth(path, reqInit);
		const payload = await parseJsonSafe(res);
		const payloadObj =
			typeof payload === 'object' && payload !== null
				? (payload as { message?: string; errors?: ValidationError[] })
				: undefined;

		if (res.ok) {
			return {
				ok: true,
				loading: false,
				data: payload as GetResponseData<P, M>
			};
		}

		return {
			ok: false,
			loading: false,
			error: {
				status: res.status,
				message: (payloadObj?.message ?? res.statusText) || 'Request failed',
				errors: payloadObj?.errors,
				body: payload
			}
		};
	} catch (e: unknown) {
		const message = e instanceof Error ? e.message : typeof e === 'string' ? e : 'Network error';
		return {
			ok: false,
			loading: false,
			error: {
				message
			}
		};
	}
}
