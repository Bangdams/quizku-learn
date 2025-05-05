<?php
// Start session if needed
session_start();
?>
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Quizku - Kelola Pengguna</title>
    <!-- Bootstrap CSS -->
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <!-- Font Awesome for icons -->
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
    <!-- Custom CSS -->
    <style>
        /* Custom Styles for Quizku Dashboard */

        /* General Styles */
        body {
            background-color: #f5f7fb;
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
        }

        /* Sidebar Styles */
        .sidebar {
            background-color: #1f2b49;
            color: #fff;
            min-height: 100vh;
            position: relative;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
            padding: 0;
        }

        .sidebar-header {
            padding: 20px 15px;
            border-bottom: 1px solid rgba(255, 255, 255, 0.1);
        }

        .sidebar-header h3 {
            margin: 0;
            font-size: 1.4rem;
            font-weight: 600;
            display: flex;
            align-items: center;
            justify-content: space-between;
        }

        .back-btn {
            font-size: 1.2rem;
            cursor: pointer;
        }

        .sidebar-nav {
            padding-top: 20px;
        }

        .sidebar-nav .nav-item {
            margin-bottom: 5px;
        }

        .sidebar-nav .nav-link {
            color: rgba(255, 255, 255, 0.7);
            padding: 10px 15px;
            border-radius: 0;
            transition: all 0.3s;
            display: flex;
            align-items: center;
        }

        .sidebar-nav .nav-link i {
            margin-right: 10px;
            width: 20px;
            text-align: center;
        }

        .sidebar-nav .nav-link:hover,
        .sidebar-nav .nav-item.active .nav-link {
            color: #ffffff;
            background-color: rgba(255, 255, 255, 0.1);
            border-left: 4px solid #ffffff;
        }

        .sidebar-nav .nav-item.active .nav-link {
            font-weight: 500;
        }

        /* Sidebar Footer with User Profile */
        .sidebar-footer {
            position: absolute;
            bottom: 0;
            width: 100%;
            padding: 15px;
            border-top: 1px solid rgba(255, 255, 255, 0.1);
        }

        .user-profile {
            display: flex;
            align-items: center;
            margin-bottom: 15px;
        }

        .user-avatar {
            width: 40px;
            height: 40px;
            border-radius: 50%;
            margin-right: 10px;
        }

        .user-info h6 {
            margin: 0;
            font-size: 0.9rem;
            font-weight: 600;
        }

        .user-info small {
            color: rgba(255, 255, 255, 0.7);
            font-size: 0.75rem;
        }

        .logout-btn {
            display: block;
            color: #fff;
            background-color: #3a4a6d;
            text-align: center;
            padding: 8px;
            border-radius: 4px;
            text-decoration: none;
            transition: background-color 0.3s;
            margin-top: 10px;
        }

        .logout-btn:hover {
            background-color: #4e5d80;
            color: #fff;
        }

        .logout-btn i {
            margin-right: 5px;
        }

        /* Main Content Area */
        .main-content {
            padding: 30px;
        }

        .page-header {
            margin-bottom: 25px;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        .page-header h2 {
            font-weight: 600;
            color: #333;
            margin-bottom: 5px;
        }

        .page-description {
            color: #6c757d;
            font-size: 0.9rem;
        }

        /* User Management Styles */
        .user-table {
            background: #fff;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
            overflow: hidden;
        }

        .user-table th {
            font-weight: 600;
            color: #333;
            background-color: #f8f9fa;
            border-bottom: 2px solid #dee2e6;
        }

        .role-badge {
            background-color: #e9ecef;
            color: #495057;
            padding: 5px 10px;
            border-radius: 50px;
            font-size: 0.8rem;
            font-weight: 500;
        }

        .action-btn {
            color: #6c757d;
            margin-right: 10px;
            transition: color 0.3s;
        }

        .edit-btn:hover {
            color: #007bff;
        }

        .delete-btn:hover {
            color: #dc3545;
        }

        .search-container {
            position: relative;
            max-width: 250px;
        }

        .search-container input {
            padding-left: 35px;
            border-radius: 50px;
            border: 1px solid #ced4da;
        }

        .search-container i {
            position: absolute;
            left: 12px;
            top: 11px;
            color: #6c757d;
        }

        .dropdown-select {
            border-radius: 50px;
            padding: 0.375rem 2rem 0.375rem 1rem;
            border: 1px solid #ced4da;
            background-position: right 12px center;
        }

        .add-user-btn {
            background-color: #1f2b49;
            border-color: #1f2b49;
            border-radius: 50px;
            padding: 0.5rem 1.5rem;
        }

        .add-user-btn:hover {
            background-color: #304269;
            border-color: #304269;
        }

        .user-avatar-sm {
            width: 30px;
            height: 30px;
            border-radius: 50%;
            margin-right: 10px;
        }

        /* Responsive adjustments */
        @media (max-width: 768px) {
            .sidebar {
                min-height: auto;
            }

            .main-content {
                padding: 15px;
            }

            .page-header {
                flex-direction: column;
                align-items: flex-start;
            }

            .add-user-btn {
                margin-top: 15px;
            }
        }
    </style>
</head>

<body>
    <div class="container-fluid p-0">
        <div class="row g-0">
            <!-- Include the sidebar navigation -->
            <?php include 'navbar.php'; ?>

            <!-- Main Content Area -->
            <div class="col-md-10 main-content">
                <div class="page-header">
                    <div>
                        <h2>Kelola Pengguna</h2>
                        <p class="page-description">Di sini kamu bisa membuat dan mengatur akun untuk Mahasiswa, Dosen, dan Admin.</p>
                    </div>
                    <button class="btn btn-primary add-user-btn">
                        <i class="fas fa-plus me-2"></i> Tambah User
                    </button>
                </div>

                <!-- Search and Filter Controls -->
                <div class="d-flex justify-content-between mb-4">
                    <div class="search-container">
                        <i class="fas fa-search"></i>
                        <input type="text" class="form-control" placeholder="Search">
                    </div>
                    <div>
                        <select class="form-select dropdown-select">
                            <option selected>All Roles</option>
                            <option>Mahasiswa</option>
                            <option>Dosen</option>
                            <option>Admin</option>
                        </select>
                    </div>
                </div>

                <!-- Users Table -->
                <div class="table-responsive user-table">
                    <table class="table table-hover">
                        <thead>
                            <tr>
                                <th>User</th>
                                <th>Role</th>
                                <th>Created</th>
                                <th>Action</th>
                            </tr>
                        </thead>
                        <tbody>
                            <?php
                            // Sample user data - in a real app, this would come from a database
                            $users = [
                                ['name' => 'Alex Smith', 'email' => 'alexsmith123@gmail.com', 'role' => 'Mahasiswa', 'created' => '2023-01-15'],
                                ['name' => 'Alex Smith', 'email' => 'alexsmith123@gmail.com', 'role' => 'Mahasiswa', 'created' => '2023-01-15'],
                                ['name' => 'Alex Smith', 'email' => 'alexsmith123@gmail.com', 'role' => 'Mahasiswa', 'created' => '2023-01-15'],
                                ['name' => 'Alex Smith', 'email' => 'alexsmith123@gmail.com', 'role' => 'Mahasiswa', 'created' => '2023-01-15'],
                                ['name' => 'Alex Smith', 'email' => 'alexsmith123@gmail.com', 'role' => 'Mahasiswa', 'created' => '2023-01-15'],
                                ['name' => 'Alex Smith', 'email' => 'alexsmith123@gmail.com', 'role' => 'Mahasiswa', 'created' => '2023-01-15'],
                                ['name' => 'Alex Smith', 'email' => 'alexsmith123@gmail.com', 'role' => 'Mahasiswa', 'created' => '2023-01-15'],
                                ['name' => 'Alex Smith', 'email' => 'alexsmith123@gmail.com', 'role' => 'Mahasiswa', 'created' => '2023-01-15'],
                                ['name' => 'Alex Smith', 'email' => 'alexsmith123@gmail.com', 'role' => 'Mahasiswa', 'created' => '2023-01-15'],
                                ['name' => 'Alex Smith', 'email' => 'alexsmith123@gmail.com', 'role' => 'Mahasiswa', 'created' => '2023-01-15']
                            ];

                            foreach ($users as $user) {
                                $formattedDate = date('Y-m-d', strtotime($user['created']));
                                echo '<tr>
                                    <td>
                                        <div class="d-flex align-items-center">
                                            <img src="https://via.placeholder.com/30" class="user-avatar-sm">
                                            <div>
                                                <div>' . $user['name'] . '</div>
                                                <small class="text-muted">' . $user['email'] . '</small>
                                            </div>
                                        </div>
                                    </td>
                                    <td><span class="role-badge">' . $user['role'] . '</span></td>
                                    <td>' . $formattedDate . '</td>
                                    <td>
                                        <a href="#" class="action-btn edit-btn"><i class="fas fa-pen"></i></a>
                                        <a href="#" class="action-btn delete-btn"><i class="fas fa-trash"></i></a>
                                    </td>
                                </tr>';
                            }
                            ?>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>

    <!-- Bootstrap JS Bundle with Popper -->
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
</body>

</html>