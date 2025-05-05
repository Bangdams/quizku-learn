<!DOCTYPE html>
<html lang="id">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login | Quizku</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <style>
        body,
        html {
            height: 100%;
            margin: 0;
        }

        .left-panel {
            background-image: url(../../Gambar/login.png);
            background-size: cover;
            position: relative;
            color: white;
        }

        .left-panel::before {
            content: '';
            position: absolute;
            top: 0;
            bottom: 0;
            left: 0;
            right: 0;
            background: rgba(0, 0, 0, 0.6);
        }

        .welcome-text {
            position: relative;
            z-index: 2;
        }

        .login-box {
            background-color: #0d1117;
            padding: 2rem;
            border-radius: 10px;
            color: white;
            width: 100%;
            max-width: 400px;
        }

        .form-control::placeholder {
            color: #999;
        }

        .login-box a {
            color: #ccc;
            font-size: 0.9rem;
            text-decoration: none;
        }

        .login-box a:hover {
            text-decoration: underline;
        }
    </style>
</head>

<body>

    <div class="container-fluid h-100">
        <div class="row h-100">
            <!-- KIRI -->
            <div class="col-md-6 d-none d-md-flex align-items-center justify-content-center left-panel">
                <div class="text-center welcome-text">
                    <h2 class="fw-bold">Selamat Datang</h2>
                    <h4 class="mt-2">di Quizku</h4>
                </div>
            </div>

            <!-- KANAN -->
            <div class="col-md-6 d-flex align-items-center justify-content-center bg-light">
                <div class="login-box">
                    <h3 class="text-center fw-bold mb-4">Login</h3>
                    <form method="POST" action="">
                        <div class="mb-3">
                            <input type="email" name="email" class="form-control" placeholder="Email" required>
                        </div>
                        <div class="mb-3">
                            <input type="password" name="password" class="form-control" placeholder="Password" required>
                        </div>
                        <div class="d-grid d-flex justify-content-center">
                            <a href="../dasboard/Dahsboard.php">
                                <button type="submit" name="login" class="btn btn-light fw-bold">Login</button>
                            </a>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>

</body>

</html>