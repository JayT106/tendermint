module.exports = {
  theme: "cosmos",
  title: "Tendermint",
  // locales: {
  //   "/": {
  //     lang: "en-US"
  //   },
  //   "/ru/": {
  //     lang: "ru"
  //   }
  // },
  base: process.env.VUEPRESS_BASE,
  themeConfig: {
    docsRepo: "tendermint/tendermint",
    editLinks: true,
    docsDir: "docs",
    label: "core",
    topbar: {
      banner: false
    },
    versions: [
      {
        "label": "v0.32",
        "key": "v0.32"
      },
      {
        "label": "v0.33",
        "key": "v0.33"
      },
      {
        "label": "v0.34",
        "key": "v0.34"
      },
      {
        "label": "master",
        "key": "master"
      }
    ],
    sidebar: {
      auto: true,
      nav: [
        {
          title: 'Resources',
          children: [
            {
              title: 'Developer Sessions',
              path: '/DEV_SESSIONS.html'
            },
            {
              title: 'RPC',
              path: 'https://docs.tendermint.com/master/rpc/',
              static: true
            },
          ]
        }
      ]
    },
    gutter: {
      title: "Help & Support",
      editLink: true,
      forum: {
        title: "Tendermint Forum",
        text: "Join the Tendermint forum to learn more",
        url: "https://forum.cosmos.network/c/tendermint",
        bg: "#0B7E0B",
        logo: "tendermint"
      },
      github: {
        title: "Found an Issue?",
        text: "Help us improve this page by suggesting edits on GitHub."
      }
    },
    footer: {
      question: {
        text: 'Chat with Tendermint developers in <a href=\'https://discord.gg/W8trcGV\' target=\'_blank\'>Discord</a> or reach out on the <a href=\'https://forum.cosmos.network/c/tendermint\' target=\'_blank\'>Tendermint Forum</a> to learn more.'
      },
      logo: "/logo-bw.svg",
      textLink: {
        text: "tendermint.com",
        url: "https://tendermint.com"
      },
      services: [
        {
          service: "medium",
          url: "https://medium.com/@tendermint"
        },
        {
          service: "twitter",
          url: "https://twitter.com/tendermint_team"
        },
        {
          service: "linkedin",
          url: "https://www.linkedin.com/company/tendermint/"
        },
        {
          service: "reddit",
          url: "https://reddit.com/r/cosmosnetwork"
        },
        {
          service: "telegram",
          url: "https://t.me/cosmosproject"
        },
        {
          service: "youtube",
          url: "https://www.youtube.com/c/CosmosProject"
        }
      ],
      smallprint:
        "The development of the Tendermint project is led primarily by Tendermint Inc., the for-profit entity which also maintains this website. Funding for this development comes primarily from the Interchain Foundation, a Swiss non-profit.",
      links: [
        {
          title: "Documentation",
          children: [
            {
              title: "Cosmos SDK",
              url: "https://cosmos.network/docs"
            },
            {
              title: "Cosmos Hub",
              url: "https://hub.cosmos.network/"
            }
          ]
        },
        {
          title: "Community",
          children: [
            {
              title: "Tendermint blog",
              url: "https://medium.com/@tendermint"
            },
            {
              title: "Forum",
              url: "https://forum.cosmos.network/c/tendermint"
            }
          ]
        },
        {
          title: "Contributing",
          children: [
            {
              title: "Contributing to the docs",
              url: "https://github.com/tendermint/tendermint"
            },
            {
              title: "Source code on GitHub",
              url: "https://github.com/tendermint/tendermint"
            },
            {
              title: "Careers at Tendermint",
              url: "https://tendermint.com/careers"
            }
          ]
        }
      ]
    },
  },
  plugins: [
    [
      "@vuepress/google-analytics",
      {
        ga: "UA-51029217-11"
      }
    ]
  ],
};