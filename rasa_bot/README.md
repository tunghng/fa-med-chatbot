# RASA Application Deployment Guide

To deploy the RASA application, follow these steps:

1. **Create your virtual environment (venv) with Python 3.8 or Python 3.9:**
    ```bash
    python3.8 -m venv venv
    ```

2. **Activate your virtual environment:**
    ```bash
    source venv/bin/activate
    ```

3. **Install RASA:**
    ```bash
    pip install rasa
    ```

4. **Training:**
    ```bash
    rasa train
    ```

5. **Chatting test:**
    ```bash
    rasa shell
    ```

6. **Integrate with Telegram:**
    - Open `credentials.yml`, replace `access_token`, `verify` (your botname), and `webhook_url` with your bot information.
    - Then run:
    ```bash
    rasa run -m models --enable-api --cors "*" --debug
    ```

    If there are any errors, try fixing with:
    ```bash
    pip install -U aiogram==2.25.2
    ```

7. **After RASA is running on your localhost:5005:**
    You can call the API through:
    - **URL:** `http://localhost:5005/webhooks/rest/webhook`
    - **Method:** `POST`
    - **Request Body:**
        ```json
        {
            "sender": "User",
            "message": "I have abdominal pain" 
        }
        ```
    The success response provides details about a successful operation.
    - **Code:** `200 OK`
    - **Content:** The response body contains the data related to the successful operation. Below is an example response for creating a post:
        ```json
        {
            "recipient_id": "User",
            "text": "/ab"
        }
        ```

**Additional Notes:**

1. Read more about RASA in the [documentation](https://rasa.com/docs/).
2. You can change `test_stories.yml` then run: `rasa test`.
3. Refine the model by changing `config.yml`.
4. Validate data by running: `rasa data validate`.
