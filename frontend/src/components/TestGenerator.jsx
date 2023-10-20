import React, { useState } from "react";
import { Formik, Field, Form, ErrorMessage } from "formik";
import axios from "axios"; // Import Axios
// Utility function
function getResultFileName(sourceFileName) {
  console.log("from getResultFileName:", sourceFileName);
  // Extract the name without the extension and append ".zip"
  const nameWithoutExtension = sourceFileName.replace(/\.[^/.]+$/, "");
  return `${nameWithoutExtension}.zip`;
}

export default function TestGenerator() {
  const initialValues = {
    sourceFile: null,
    numFiles: 1,
    examTitle: "",
    beforeTest: "",
    merge: true,
    newPage: true,
  };
  const [selectedFile, setSelectedFile] = useState(null);
  const onSubmit = async (values, actions) => {
    try {
      console.log(selectedFile);
      console.log("values: ", values.numFiles.toString());
      const formData = new FormData();
      formData.append("sourceFile", selectedFile);
      formData.append("resultFile", getResultFileName(selectedFile.name));
      formData.append("numFiles", values.numFiles.toString());
      formData.append("examTitle", values.examTitle);
      formData.append("beforeTest", values.beforeTest);
      formData.append("merge", values.merge);
      formData.append("newPage", values.newPage);
      console.log("form below:");
      console.log("form data: ", formData);
      // Post Request to API
      const response = await axios.post("/generate-test", formData);

      // Handle the Api response
      console.log(response.data);
    } catch (error) {
      // Handle errors
      console.error(error);
    }
  };
  return (
    <div>
      <h2>Test Generator</h2>
      <Formik initialValues={initialValues} onSubmit={onSubmit}>
        <Form>
          <div>
            <label htmlFor="sourceFile">Select a Text File</label>
            <input
              type="file"
              id="sourceFile"
              name="sourceFile"
              onChange={(event) =>
                setSelectedFile(event.currentTarget.files[0])
              }
            />
            <ErrorMessage name="sourceFile" component="div" />
          </div>

          <div>
            <label htmlFor="numFiles">Number of Files</label>
            <Field type="number" id="numFiles" name="numFiles" />
            <ErrorMessage name="numFiles" component="div" />
          </div>

          <div>
            <label htmlFor="examTitle">ExamTitle:</label>
            <Field type="text" id="examTitle" name="examTitle" />
            <ErrorMessage name="examTitle" component="div" />
          </div>

          <div>
            <label htmlFor="beforeTest">Before Test</label>
            <Field type="text" id="beforeTest" name="beforeTest" />
            <ErrorMessage name="beforeTest" component="div"></ErrorMessage>
          </div>

          <div>
            <label>
              Merge:
              <Field type="checkbox" name="merge" />
            </label>
          </div>

          <div>
            <label>
              New Page:
              <Field type="checkbox" name="newPage" />
            </label>
          </div>

          <button type="submit">GO</button>
        </Form>
      </Formik>
    </div>
  );
}
