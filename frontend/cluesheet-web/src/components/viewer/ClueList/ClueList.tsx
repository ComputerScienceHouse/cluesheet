"use client"
import { Clue } from "@/lib/types";
import styles from "./page.module.scss";

/*const transform = (data: Clue, parent: Clue) =>  {
    return Object.keys(data).map((key) => {
        const value = data[key];
        const node = {
            label: key,
            checked: false,
            childrenNodes: [],
            parent: parent,
        };

        if (typeof value === "boolean") {
            node.checked = value;
        } else {
            const children = transform(value, node);
            node.childrenNodes = children;
            if (children.every((node) => node.checked)) {
                node.checked = true;
            }
        }

        return node;
    });
}*/

export function ClueItem() {}

interface ClueItemProps {
  ancestors: Array<Clue>;
  clues: Array<Clue>;
  edit: boolean;
}

function handleClueChecked() {
  console.log("Checked!");
}

function getInputField(clue: Clue) {
  if (clue.rule?.key === "stacks") {
    return (<input type="number"/>);
  }

  return (
    <input
      type="checkbox"
      className={styles.customCheckboxInput}
      onChange={handleClueChecked}
      checked={clue.checked}
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
                getInputField(clue)
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
