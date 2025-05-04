export const mockCluesheetUUID = "3a4d3998-71cb-42ca-942c-202db9e465c4";

export const mockCluesheet = {
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
      rule: null,
      tags: [],
      children: [
        {
          id: "my-sick-and-poggers-uuid-2",
          limit: 1,
          completions: 0,
          description: "No utensils",
          points: "+1",
          rule: null,
          tags: [],
          children: [
            {
              id: "my-sick-and-poggers-uuid-3",
              limit: 1,
              completions: 0,
              description: "With a straw",
              points: "+1",
              rule: null,
              tags: [],
              children: [
                {
                  id: "my-sick-and-poggers-uuid-4",
                  limit: 1,
                  completions: 0,
                  description: "Is a straw a utensil?",
                  points: "+1",

                  rule: null,
                  tags: [],
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

          rule: null,
          tags: [],
          children: [],
        },
      ],
    },
    {
      id: "my-sick-and-poggers-uuid-6",
      limit: 0,
      completions: 3,
      description: "Fix something broken (stacks)",
      points: "+10",

      rule: null,
      tags: [],
      children: [],
    },
  ],
};
