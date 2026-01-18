import { API_BASE_URL } from '$lib/constants';
import { FirebaseAuthentication } from '@capacitor-firebase/authentication';

export async function fetchWithAuth(url: string, options: RequestInit = {}) {
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
