# rex
Text utility functions that use regexp searches -- I use mainly for JSON cleanup

Useful to compress the json.Marshall output from something like:

    {
      "title": "Kick The Foobar",
      "artist": "Foo & The Bars (f/ Boo)",
      "number": 5,
      "trk-stt": "16:20'180",
      "trk-end": "24:43'899!48",
      "trk-len": "8:23.800",
      "aud-stt": "16:20'414",
      "aud-end": "24:39'163!48",
      "aud-len": "8:18.722",
      "silences": [
        {
          "stt": "22:13'395!32",
          "end": "22:13'742!20"
        },
        {
          "stt": "22:15'102!08",
          "end": "22:16'4!42"
        },
        {
          "stt": "23:37'217!22",
          "end": "23:38'193!28"
        }
      ]
    },

Can be compressed to a more human friendly looking:

    { "title":        "Kick The Foobar",
      "artist":       "Foo & The Bars (f/ Boo)",
      "number": 5,
      "trk-stt": "16:20'180", "trk-end": "24:43'899!48", "trk-len": "8:23.800",
      "aud-stt": "16:20'414", "aud-end": "24:39'163!48", "aud-len": "8:18.722",
      "silences": [
        { "stt": "22:13'395!32", "end": "22:13'742!20" },
        { "stt": "22:15'102!08", "end": "22:16'4!42" },
        { "stt": "23:37'217!22", "end": "23:38'193!28" }
      ]
    },
