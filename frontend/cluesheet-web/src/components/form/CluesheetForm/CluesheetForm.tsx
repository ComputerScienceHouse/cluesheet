"use client";

import { cluesheetBackendEndpoint } from "@/lib/endpoint";
import styles from "./page.module.scss";
import { Container } from "reactstrap";
import FlipNumbers from "react-flip-numbers";
import PointCounter from "@/components/viewer/counter/PointCounter";
import { Clue } from "@/lib/types";
import { ClueList } from "@/components/viewer/ClueList/ClueList";
import { mockUserCluesheet } from "../../../../tests/lib/data";
import { useEffect, useState } from "react";

interface CluesheetFormProps {
  cluesheet_id: string;
  keycloak_uid: string;
}

export default function CluesheetForm({
  cluesheet_id,
  keycloak_uid,
}: CluesheetFormProps) {
  //const { cluesheet_id } = await params;
  //const cluesheet = await fetch(`${cluesheetBackendEndpoint}/api/v1/cluesheet/${cluesheet_id}/user/${keycloak_uid}`);
  //console.log(`Got cluesheet object: ${await cluesheet.json()}`);

  const cluesheet = mockUserCluesheet;
  const [clues, setClues] = useState(cluesheet.clues);
  const [score, setScore] = useState(cluesheet.user_points);

  const handleCheckboxChange = (clueId: string, completionValue: number = NaN) => {
    console.log(`changed ${clueId}`);

    function updateClue(clue: Clue) {
      console.log(clue);
      return clue;
    }

    function updateClues(clues: Array<Clue>): Array<Clue> {
      const updatedClues = clues.map((clue: Clue) => {
        clue.children = updateClues(clue.children);
        clue = updateClue(clue);
        return clue;
      });
      return updatedClues;
    }

    updateClues(clues);


    /*
    setClues((clues) => {
      const updatedClues = clues.map((clue) => {
        if (clue.id === clueId) {
          console.log(`completions = ${clue.completions}`);
          // Set the number of completions
          if (!isNaN(completionValue)) {
            return { ...clue, completions: completionValue }
          }
          // If the checkbox is not checked, then check it
          if (clue.completions === 0) {
            return { ...clue, completions: 1 };
          }
          // If the checkbox is checked, then un-check it
          return { ...clue, completions: 0 };
        }
        return clue;
      });
      return updatedClues;
    });
    */
  };

  function convertToInteger(input: string): number {
    // Trim the input to remove any leading or trailing whitespace
    const trimmedInput = input.trim();

    // Use a regular expression to extract the numeric part
    const match = trimmedInput.match(/[-+]?\d+/);

    // If a match is found, convert it to an integer; otherwise, return NaN
    return match ? parseInt(match[0], 10) : NaN;
  }

  const calculatePoints = () => {
    let total = 0;
    clues.map((clue) => {
      const cluePoints = convertToInteger(clue.points);
      total += cluePoints * clue.completions;
    });
    return total;
  };

  useEffect(() => {
    console.log(clues);
    setScore(calculatePoints());
  }, [clues]);

  return (
    <>
      <div className={styles.header}>
        <h1>{cluesheet.title}</h1>
        <div className={styles.pointCounter}>
          <PointCounter points={score} />
          <h3>Points</h3>
        </div>
      </div>
      <div className={styles.clueList}>
        {cluesheet.clues.length > 0 && (
          <ClueList
            ancestors={[]}
            clues={clues}
            edit={true}
            handleCheckboxChange={handleCheckboxChange}
          />
        )}
      </div>
    </>
  );
}
