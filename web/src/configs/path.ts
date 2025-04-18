/**
 * The base path of the application. Uses '/' in development and a placeholder in production, which will be replaced at
 * runtime with the actual base path.
 */
export const basename = import.meta.env.DEV ? '/' : '/__KOALA_BASE_PATH__';

/** The base path for API endpoints. Appends '/api' to the application's base path. */
export const apiBase = basename === '/' ? '/api' : basename + '/api';
