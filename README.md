This tiny app can be run as part of post-processing in Paperless-NGX.

It will replace the content extracted by crappy tesseract-ocr with much better output provided by Google DocumentAI API.

The app requires a few parameters (can be set via env as well):

GOOGLE_APPLICATION_CREDENTIALS=path to your service JSON file form document API console
DOCUMENTAI_PROJECT_ID=your project id
DOCUMENTAI_LOCATION=eu/us
DOCUMENTAI_PROCESSOR_ID=the ocr processor created in document API console
DOCUMENTAI_PROCESSOR_VERSION=rc
PAPERLESS_TOKEN=your user's token in paperless
PAPERLESS_ENDPOINT=http://localhost:8000

To use this with docker-compose'd paperless-ngx:

```
    volumes:
      .......
      - ./postprocess:/postprocess
    environment:
      .......
      PAPERLESS_POST_CONSUME_SCRIPT: /postprocess/run.sh
```
