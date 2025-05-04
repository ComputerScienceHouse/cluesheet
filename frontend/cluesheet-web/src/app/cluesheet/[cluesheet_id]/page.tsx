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
  return (
    <>
      <main>
      </main>
    </>
  );
}
