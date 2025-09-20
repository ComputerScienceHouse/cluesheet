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

  const handleCheckboxChange = (clueId: string) => {
    console.log(`changed ${clueId}`);
    setClues((clues) => {
      const updatedClues = clues.map((clue) => {
        if (clue.id === clueId) {
          return { ...clue, checked: !clue.checked };
        }
        return clue;
      });
      return updatedClues;
    });
  };

  useEffect(() => {
    console.log(clues);
  }, [clues]);

  return (
    <>
      <div className={styles.header}>
        <h1>{cluesheet.title}</h1>
        <div className={styles.pointCounter}>
          <PointCounter points={cluesheet.user_points} />
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
