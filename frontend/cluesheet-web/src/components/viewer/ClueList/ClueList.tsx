"use client";
import { Clue } from "@/lib/types";
import styles from "./page.module.scss";
import InputField from "../InputField/InputField";
import { LineItem } from "../LineItem/LineItem";

export function ClueItem() {}

interface ClueItemProps {
  ancestors: Array<Clue>;
  clues: Array<Clue>;
  edit: boolean;
}

export function ClueList({ ancestors, clues, edit = false }: ClueItemProps) {
  const ancestorTags: string[] = ancestors.flatMap((ancestor) => ancestor.tags);

  return (
    <ul>
      {clues.map((clue: Clue, index) => (
        <li key={clue.id} className={styles.clueBody}>
          <LineItem ancestors={ancestors} clue={clue} />
          {clue.children.length > 0 && (
            <ClueList
              ancestors={[...ancestors, clue]}
              clues={clue.children}
              edit={edit}
            />
          )}
        </li>
      ))}
    </ul>
  );
}
