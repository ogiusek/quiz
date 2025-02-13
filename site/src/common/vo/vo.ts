/**
 * VO standed for value object.
 * but this is just validation interface with 2 methods one for listing all errors and second for listing errors which occured with current values
 */
export interface VO {
  Valid(): Error[]
  Errors(): Error[]
} 