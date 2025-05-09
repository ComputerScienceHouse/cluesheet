"use client"
import { Clue } from "@/lib/types";
import styles from "./page.module.scss";
import { useState } from "react";

export function ClueItem() {}

interface ClueItemProps {
  ancestors: Array<Clue>;
  clues: Array<Clue>;
  edit: boolean;
}

function getInputField(clue: Clue, ancestors: Array<Clue>) {
  if (clue.rule?.key === "stacks") {
    return (<input type="number"/>);
  }

  // state
  const [checked, setChecked] = useState(clue.checked ?? false);
  
  // checkbox click handler
  function handleClick(e) {
    console.log(`${clue.id} checked`);
    setChecked(!checked);
  };

  return (
    <input
      type="checkbox"
      value={clue.id}
      className={styles.customCheckboxInput}
      onChange={handleClick}
      checked={checked}
    />
  );
}

export function ClueList({ ancestors, clues, edit = false }: ClueItemProps) {
  const ancestorTags: string[] = ancestors.flatMap((ancestor) => ancestor.tags);

  return (
    <ul>
      {clues.length > 0 &&
        clues.map((clue: Clue, index) => (
          <li key={clue.id} className={styles.clueBody}>
            <div className={styles.lineItem}>
              { edit &&
                getInputField(clue, ancestors)
              }
              <div className={styles.points}>{clue.points} {clue.rule && `(${clue.rule.key})`}</div>
              <div className={styles.description}>{clue.description}</div>
              <div className={styles.tags}>
                {clue.tags
                  .filter((tag) => !ancestorTags.includes(tag))
                  .map((tag: string, index) => (
                    <p className={styles.tag} key={index}>
                      {tag}
                    </p>
                  ))}
              </div>
            </div>

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
