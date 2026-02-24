export const mockCluesheetUUID = "3a4d3998-71cb-42ca-942c-202db9e465c4";

export const mockCluesheet = {
  id: "my-epic-cluesheet-uuid",
  title: "Yet-another-Rip n' Dip-a-thon 2192",
  clues: [
    {
      id: "my-sick-and-poggers-uuid",
      description: "Linux",
      points: "+1",
      rule: null,
      tags: ["Operating Systems", "Computers", "FOSS"],
      children: [
        {
          id: "my-sick-and-poggers-uuid-100",
          description: "Arch",
          points: "+2",
          rule: null,
          tags: [],
          children: [],
        },
        {
          id: "my-sick-and-poggers-uuid-101",
          description: "Install Linux",
          points: "+1",
          rule: {
            key: "stacks",
            description:
              "You can complete this rule as many times as you want.",
          },
          tags: [],
          children: [
            {
              id: "my-sick-and-poggers-uuid-100",
              description: "Not your computer",
              points: "+4",
              rule: {
                key: "stacks",
                description:
                  "You can complete this rule as many times as you want.",
              },
              tags: ["Operating Systems", "Computers", "FOSS"],
              children: [],
            },
          ],
        },
      ],
    },
    {
      id: "my-sick-and-poggers-uuid-6",
      description: "Fix something broken",
      points: "+10",
      rule: {
        key: "stacks",
        description: "You can complete this clue as many times as you want.",
      },
      tags: ["Opcommathon 2023"],
      children: [],
    },
  ],
};

export const mockUserCluesheet = {
  id: "my-epic-cluesheet-uuid",
  title: "Opcommathon 2069 Cluesheet '>w<'",
  user_points: 69420,
  user_additional_score: ["-Your Bones", "+Adam Neulight's Car"],
  /*point_unit: "Kubernetes Clusters",*/
  clues: [
    {
      id: "my-sick-and-poggers-uuid",
      completions: 0,
      description: "Eat a whole can of beans",
      points: "+5",
      rule: null,
      tags: ["beans"],
      children: [
        {
          id: "my-sick-and-poggers-uuid-2",
          completions: 0,
          description: "No utensils",
          points: "+1",
          rule: null,
          tags: [],
          children: [
            {
              id: "my-sick-and-poggers-uuid-3",
              completions: 0,
              description: "With a straw",
              points: "+1",
              rule: null,
              tags: [],
              children: [
                {
                  id: "my-sick-and-poggers-uuid-4",
                  completions: 0,
                  description: "Is a straw a utensil?",
                  points: "+1",
                  rule: null,
                  tags: ["beans"],
                  children: [],
                },
              ],
            },
          ],
        },
        {
          id: "my-sick-and-poggers-uuid-5",
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
      completions: 3,
      description: "Fix something broken",
      points: "+10",
      rule: {
        key: "stacks",
        description: "You can complete this clue as many times as you want.",
      },
      tags: [],
      children: [],
    },
    {
      id: "my-sick-and-poggers-uuid-26",
      completions: 0,
      description: "Die",
      points: "+100000",
      rule: null,
      tags: ["death"],
      children: [],
    },
  ],
};
