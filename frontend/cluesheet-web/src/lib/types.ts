export interface Clue {
  id: string;
  description: string;
  points: string;
  rule: Rule | null;
  tags: Array<string>;
  children: Array<Clue>;
};

export interface UserClue extends Clue {
  completions: number; // Number of times this has been completed
}

export type Rule = {
  key: string;
  description: string;
};
