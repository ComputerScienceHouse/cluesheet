import { cluesheetBackendEndpoint } from "@/lib/endpoint";

export const metadata = {
  title: "View cluesheet",
  description: "View an mf cluesheet kerchoo",
};

export default async function CluesheetViewer({
  params,
}: {
  params: Promise<{ cluesheet_id: string }>;
}) {
  const { cluesheet_id } = await params;

  const cluesheet = await fetch(`${cluesheetBackendEndpoint}/${cluesheet_id}`);

  console.log(`Got cluesheet object: ${await cluesheet.json()}`);
  
  return (
    <>
      <main>
      </main>
    </>
  );
}
