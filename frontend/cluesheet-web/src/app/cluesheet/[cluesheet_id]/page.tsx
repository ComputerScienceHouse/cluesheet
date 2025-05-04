import { cluesheetBackendEndpoint } from "@/lib/endpoint";
import styles from "./page.module.scss";
import { Container } from "reactstrap";
import FlipNumbers from "react-flip-numbers";
import PointCounter from "@/components/viewer/counter/PointCounter";
import { Clue } from "@/lib/types";
import ClueGroup from "@/components/viewer/cluegroup/ClueGroup";

export const metadata = {
  title: "View cluesheet",
  description: "View an mf cluesheet kerchoo",
};

export default async function CluesheetViewer({
  params,
}: {
  params: Promise<{ cluesheet_id: string }>;
}) {
  //const { cluesheet_id } = await params;
  //const cluesheet = await fetch(`${cluesheetBackendEndpoint}/api/v1/cluesheet/${cluesheet_id}`);
  //console.log(`Got cluesheet object: ${await cluesheet.json()}`);

  const cluesheet = {
    id: "my-epic-cluesheet-uuid",
    title: "Opcommathon 2069 Cluesheet '>w<'",
    user_points: 69420,
    user_additional_score: ["-Your Bones", "+Adam Neulight's Car"],
    /*point_unit: "Kubernetes Clusters",*/
    clues: [
      {
        id: "my-sick-and-poggers-uuid",
        limit: 1,
        completions: 1,
        description: "Eat a whole can of beans",
        points: "+5",
        children: [
          {
            id: "my-sick-and-poggers-uuid-2",
            limit: 1,
            completions: 0,
            description: "No utensils",
            points: "+1",
            children: [
              {
                id: "my-sick-and-poggers-uuid-3",
                limit: 1,
                completions: 0,
                description: "With a straw",
                points: "+1",
                children: [
                  {
                    id: "my-sick-and-poggers-uuid-4",
                    limit: 1,
                    completions: 0,
                    description: "Is a straw a utensil?",
                    points: "+1",
                    children: [],
                  },
                ],
              },
            ],
          },
          {
            id: "my-sick-and-poggers-uuid-5",
            limit: 1,
            completions: 0,
            description: "On a bike",
            points: "+1",
            children: [],
          }
        ],
      },
      {
            id: "my-sick-and-poggers-uuid-6",
            limit: 0,
            completions: 3,
            description: "Fix something broken (stacks)",
            points: "+10",
            children: [],
      },
    ],
  };

  return (
    <>
      <main>
        <Container>
          <div className={styles.header}>
            <h1>{cluesheet.title}</h1>
            <div className={styles.pointCounter}>
              <PointCounter points={cluesheet.user_points} />
              <h3>Points</h3>
            </div>
          </div>
          <div className={styles.clueList}>
            {cluesheet.clues.length > 0 &&
              cluesheet.clues.map((clue: Clue, index) => (
                <ClueGroup key={`${clue.id}-${index}`} clue={clue} />
              ))}
          </div>
        </Container>
      </main>
    </>
  );
}
