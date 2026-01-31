import type { Plant } from '$lib/types/api';

interface PlantsStore {
	plants: Plant[];
	loading: boolean;
	error: string | null;

	setPlants(plants: Plant[]): void;
	setLoading(loading: boolean): void;
	setError(error: string | null): void;
	reset(): void;
}

export function createPlantsStore(): PlantsStore {
	let plants = $state<Plant[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	return {
		get plants() {
			return plants;
		},
		get loading() {
			return loading;
		},
		get error() {
			return error;
		},
		setPlants(newPlants: Plant[]) {
			plants = newPlants;
		},
		setLoading(newLoading: boolean) {
			loading = newLoading;
		},
		setError(newError: string | null) {
			error = newError;
		},
		reset() {
			plants = [];
			loading = true;
			error = null;
		}
	};
}

// Global plants store singleton
let plantsStore: PlantsStore | null = null;

export function getPlantsStore(): PlantsStore {
	if (!plantsStore) {
		plantsStore = createPlantsStore();
	}
	return plantsStore;
}
