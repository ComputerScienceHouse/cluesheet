export type Clue = {
  id: string;
  limit: number;
  completions: number;
  description: string;
  points: string;
  children: Array<Clue>;
};
