export interface Clue {
  id: string;
  completions: number; // Number of times this has been completed
  description: string;
  points: string;
  rule: Rule | null;
  tags: Array<string>;
  children: Array<Clue>;
}

export type Rule = {
  key: string;
  description: string;
};
