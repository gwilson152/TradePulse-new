import { dasTraderSchema } from './dasTrader';
import { propReportsSchema } from './propReports';
import type { PlatformSchema } from '$lib/types/import';

export const platforms: PlatformSchema[] = [
	dasTraderSchema,
	propReportsSchema
	// Add more platforms here
];

export function getPlatformById(id: string): PlatformSchema | undefined {
	return platforms.find((p) => p.id === id);
}

export function getPlatformByName(name: string): PlatformSchema | undefined {
	return platforms.find((p) => p.name.toLowerCase() === name.toLowerCase());
}

export { dasTraderSchema, propReportsSchema };
