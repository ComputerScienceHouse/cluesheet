export interface Clue {
  id: string;
  description: string;
  points: string;
  rule: Rule | null;
  tags: Array<string>;
  children: Array<Clue>;

  completions: number; // Number of times this has been completed

  locked: boolean;
  checked: boolean;
}

/*
export interface UserClue extends Clue {
  completions: number; // Number of times this has been completed
}

export interface ClueListUserClue extends UserClue {
  parent: ClueListUserClue;
  locked: boolean;
  checked: boolean;
}*/

export type Rule = {
  key: string;
  description: string;
};
